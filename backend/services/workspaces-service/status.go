package workspaceservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	messagebroker "github.com/wafi11/workspaces/pkg/message-broker"
)


func (r *Repository) lockWorkspaceForUpdate(ctx context.Context, tx *sql.Tx, workspaceId string) (*workspaceRow, error) {
	
    var res workspaceRow
    err := tx.QueryRowContext(ctx, `
        SELECT w.user_id, w.name, w.status,
            wr.limit_ram_mb, wr.limit_cpu_cores,
            wr.ram_mb_req, wr.cpu_cores_req
        FROM workspaces w
        JOIN workspace_resources wr ON wr.workspace_id = w.id
        WHERE w.id = $1
        FOR UPDATE`, workspaceId,
    ).Scan(&res.UserId, &res.Name, &res.CurrStatus,
        &res.LimitRAM, &res.LimitCPU, &res.ReqRAM, &res.ReqCPU)
    if err != nil {
        return nil, fmt.Errorf("workspace tidak ditemukan: %w", err)
    }
    return &res, nil
}

func (r *Repository) UpdateWorkspaceStatus(ctx context.Context, workspaceId string, status string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := r.lockWorkspaceForUpdate(ctx, tx, workspaceId)
	if err != nil {
		return err
	}

	if string(res.CurrStatus) == status {
		return fmt.Errorf("workspace sudah dalam status %s", status)
	}

	// Validate transition and call the appropriate handler.
	// Each handler is responsible for: quota, session, event publish, scheduler.
	switch status {
	case "running":
		if res.CurrStatus != "stopped" && res.CurrStatus != "pending" {
			return fmt.Errorf("start hanya bisa dari stopped atau pending, status saat ini: %s", res.CurrStatus)
		}
		if err := r.handleStart(ctx, tx, workspaceId, res); err != nil {
			return err
		}

	case "paused":
		if res.CurrStatus != "running" {
			log.Printf("paused calling")
			return fmt.Errorf("pause hanya bisa saat running, status saat ini: %s", res.CurrStatus)
		}
		if err := r.handlePause(ctx, tx, workspaceId, res); err != nil {
			return err
		}

	case "resumed":
		if res.CurrStatus != "paused" {
			return fmt.Errorf("resume hanya bisa dari paused, status saat ini: %s", res.CurrStatus)
		}
		if err := r.handleResume(ctx, tx, workspaceId, res); err != nil {
			return err
		}

	// case "stopped":
	// 	if res.CurrStatus != "running" {
	// 		return fmt.Errorf("stop hanya bisa saat running, status saat ini: %s", res.CurrStatus)
	// 	}
	// 	if err := r.handleStop(ctx, tx, workspaceId, res); err != nil {
	// 		return err
	// 	}

	default:
		return fmt.Errorf("status tidak valid: %s", status)
	}

	// "resumed" di DB tetap disimpan sebagai "running"
	dbStatus := status
	if status == "resumed" {
		dbStatus = "running"
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE workspaces SET status = $1, updated_at = NOW() WHERE id = $2`,
		dbStatus, workspaceId,
	)
	if err != nil {
		return fmt.Errorf("gagal update status workspace: %w", err)
	}

	return tx.Commit()
}

// AutoStopWorkspace dipanggil oleh asynq worker ketika expires_at tercapai.
// Berbeda dari handleStop (manual) — tidak ada cooldown next_start_at
// karena ini bukan user yang stop, tapi sistem.
func (r *Repository) AutoStopWorkspace(ctx context.Context, workspaceId string) error {
	log.Printf("[autostop] mulai proses workspace %s", workspaceId)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[autostop] gagal begin tx workspace %s: %s", workspaceId, err.Error())
		return err
	}
	defer tx.Rollback()

	res, err := r.lockWorkspaceForUpdate(ctx, tx, workspaceId)
	if err != nil {
		log.Printf("[autostop] gagal lock workspace %s: %s", workspaceId, err.Error())
		return err
	}
	log.Printf("[autostop] workspace %s status=%s cpu=%v ram=%v user=%s", workspaceId, res.CurrStatus, res.ReqCPU, res.ReqRAM, res.UserId)

	if res.CurrStatus == "stopped" || res.CurrStatus == "paused" {
		log.Printf("[autostop] workspace %s sudah dalam status %s, skip", workspaceId, res.CurrStatus)
		return tx.Commit()
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE user_quotas
		SET used_cpu_cores = GREATEST(used_cpu_cores - $1, 0),
		    used_ram_mb    = GREATEST(used_ram_mb - $2, 0)
		WHERE user_id = $3`,
		res.ReqCPU, res.ReqRAM, res.UserId,
	)
	if err != nil {
		log.Printf("[autostop] gagal release quota user %s: %s", res.UserId, err.Error())
		return fmt.Errorf("autostop: gagal release quota: %w", err)
	}
	log.Printf("[autostop] quota dirilis user %s cpu=%v ram=%v", res.UserId, res.ReqCPU, res.ReqRAM)

	
	var sessionID string
	err = tx.QueryRowContext(ctx, `
		SELECT id FROM workspace_sessions
		WHERE workspace_id = $1
		AND status IN ('running', 'paused')
		ORDER BY created_at DESC
		LIMIT 1
		`, workspaceId).Scan(&sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[autostop] tidak ada session aktif untuk workspace %s", workspaceId)
			return fmt.Errorf("autostop: tidak ada session aktif untuk workspace %s", workspaceId)
		}
		log.Printf("[autostop] gagal ambil session workspace %s: %s", workspaceId, err.Error())
		return fmt.Errorf("autostop: gagal ambil session: %w", err)
	}
	log.Printf("[autostop] session aktif ditemukan id=%s", sessionID)
	nextStartAt := time.Now().UTC().Add(cooldown)

	_, err = tx.ExecContext(ctx, `
		UPDATE workspace_sessions
		SET status        = 'stopped',
			next_start_at = NOW() + $1::interval,
			stopped_at    = NOW(),
			paused_at     = NULL,
			updated_at    = NOW()
		WHERE id = $2
	`, fmt.Sprintf("%d minutes", int(cooldown)), sessionID)
	log.Printf("[autostop] session %s diupdate ke stopped, next_start_at=%s", sessionID, nextStartAt.Format(time.RFC3339))

	_, err = tx.ExecContext(ctx,
		`UPDATE workspaces SET status = 'stopped', updated_at = NOW() WHERE id = $1`,
		workspaceId,
	)
	if err != nil {
		log.Printf("[autostop] gagal update workspace %s: %s", workspaceId, err.Error())
		return fmt.Errorf("autostop: gagal update workspace: %w", err)
	}
	log.Printf("[autostop] workspace %s diupdate ke stopped", workspaceId)

	r.publishStop(ctx, workspaceId, res)
	log.Printf("[autostop] stop event dipublish untuk workspace %s", workspaceId)

	if err := tx.Commit(); err != nil {
		log.Printf("[autostop] gagal commit tx workspace %s: %s", workspaceId, err.Error())
		return err
	}

	log.Printf("[consumer] sending to userID=%s clients=%v", res.UserId, r.hub)
	r.hub.SendToUser(res.UserId, map[string]string{
		"type":         fmt.Sprintf("workspace.%s","stopped"),
		"workspace_id": workspaceId,
		"status":       "stopped",
	})

	log.Printf("[autostop] selesai workspace %s berhasil distop", workspaceId)
	return nil
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

// handleStart: stopped → running
// Cek cooldown, cek quota, buat session baru, publish event, jadwalkan auto stop.
func (r *Repository) handleStart(ctx context.Context, tx *sql.Tx, workspaceId string, res *workspaceRow) error {
	// Cek cooldown (next_start_at dari session sebelumnya)
	acc, err := r.CanStartWorkspace(ctx, workspaceId, tx)
	if err != nil {
		return err
	}
	if !acc {
		return fmt.Errorf("workspace dalam cooldown, coba beberapa menit lagi")
	}

	// Cek dan klaim quota
	if err := r.claimQuota(ctx, tx, res); err != nil {
		return err
	}

	// Buat session baru
	_, err = tx.ExecContext(ctx, `
		INSERT INTO workspace_sessions (
			workspace_id, user_id, status,
			started_at, expires_at,
			created_at, updated_at
		) VALUES ($1, $2, 'running', NOW(), NOW() + $3::interval, NOW(), NOW())
	`, workspaceId, res.UserId, fmt.Sprintf("%d minutes", int(cooldown)))
	if err != nil {
		return fmt.Errorf("gagal buat session: %w", err)
	}

	if res.CurrStatus == "stopped" {
		r.publishStart(ctx, workspaceId, res)
	}
	r.scheduleAutoStop(workspaceId, cooldown)

	return nil
}

// handlePause: running → paused
// Cancel timer, release quota, tandai paused_at (expires_at TIDAK diubah supaya sisa waktu terjaga).
func (r *Repository) handlePause(ctx context.Context, tx *sql.Tx, workspaceId string, res *workspaceRow) error {
	// Cancel timer duluan sebelum apapun
	messagebroker.CancelAutoStop(workspaceId, r.redis.Inceptor)

	// Release quota sementara
	if err := r.releaseQuota(ctx, tx, res); err != nil {
		return fmt.Errorf("gagal release quota saat pause: %w", err)
	}

	// Set paused_at — expires_at SENGAJA tidak diubah, dipakai handleResume untuk hitung sisa waktu
	_, err := tx.ExecContext(ctx, `
		UPDATE workspace_sessions
		SET status     = 'paused',
		    paused_at  = NOW(),
		    updated_at = NOW()
		WHERE workspace_id = $1 AND status = 'running'
	`, workspaceId)
	if err != nil {
		return fmt.Errorf("gagal update session saat pause: %w", err)
	}

	// Publish stop ke operator (pod di-suspend/scale-down)
	r.publishStop(ctx, workspaceId, res)

	return nil
}

// handleResume: paused → running
// Hitung sisa waktu dari (expires_at - paused_at), klaim quota kembali, jadwalkan timer baru.
func (r *Repository) handleResume(ctx context.Context, tx *sql.Tx, workspaceId string, res *workspaceRow) error {
	// Cek dan klaim quota kembali
	if err := r.claimQuota(ctx, tx, res); err != nil {
		log.Printf("failed to resume workspaces : %s",err.Error())
		return err
	}

	// Ambil sisa waktu dari session yang paused
	var sessionId string
	var pausedAt, expiresAt time.Time
	err := tx.QueryRowContext(ctx, `
		SELECT id, paused_at, expires_at
		FROM workspace_sessions
		WHERE workspace_id = $1 AND status = 'paused'
		ORDER BY created_at DESC
		LIMIT 1
	`, workspaceId).Scan(&sessionId, &pausedAt, &expiresAt)
	if err != nil {
		log.Printf("failed to resumed workspaces : %s",err.Error())
		return fmt.Errorf("session paused tidak ditemukan: %w", err)
	}

	// Sisa waktu = waktu yang belum terpakai sebelum di-pause
	remainingTime := expiresAt.Sub(pausedAt)
	if remainingTime <= 0 {
		// Waktu sudah habis saat di-pause, beri grace period minimal
		remainingTime = 1 * time.Minute
	}

	newExpiresAt := time.Now().UTC().Add(remainingTime)

	_, err = tx.ExecContext(ctx, `
		UPDATE workspace_sessions
		SET status     = 'running',
		    paused_at  = NULL,
		    expires_at = $1,
		    updated_at = NOW()
		WHERE id = $2
	`, newExpiresAt, sessionId)
	if err != nil {
		log.Printf("failed to resumed workspaces : %s",err.Error())
		return fmt.Errorf("gagal resume session: %w", err)
	}

	r.publishStart(ctx, workspaceId, res)
	r.scheduleAutoStop(workspaceId, remainingTime)

	return nil
}

// ---------------------------------------------------------------------------
// Quota helpers
// ---------------------------------------------------------------------------

func (r *Repository) claimQuota(ctx context.Context, tx *sql.Tx, res *workspaceRow) error {
	var maxCPU, usedCPU float64
	var maxRAM, usedRAM int
	err := tx.QueryRowContext(ctx, `
		SELECT max_cpu_cores, used_cpu_cores, max_ram_mb, used_ram_mb
		FROM user_quotas WHERE user_id = $1 FOR UPDATE`, res.UserId,
	).Scan(&maxCPU, &usedCPU, &maxRAM, &usedRAM)
	if err != nil {
		return fmt.Errorf("gagal ambil quota: %w", err)
	}
	if usedCPU+res.ReqCPU > maxCPU {
		return fmt.Errorf("quota CPU tidak cukup, matikan workspace lain dulu")
	}
	if usedRAM+res.ReqRAM > maxRAM {
		return fmt.Errorf("quota RAM tidak cukup, matikan workspace lain dulu")
	}
	_, err = tx.ExecContext(ctx, `
		UPDATE user_quotas
		SET used_cpu_cores = used_cpu_cores + $1,
		    used_ram_mb    = used_ram_mb + $2
		WHERE user_id = $3`,
		res.ReqCPU, res.ReqRAM, res.UserId,
	)
	if err != nil {
		return fmt.Errorf("gagal update quota: %w", err)
	}
	return nil
}

func (r *Repository) releaseQuota(ctx context.Context, tx *sql.Tx, res *workspaceRow) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE user_quotas
		SET used_cpu_cores = GREATEST(used_cpu_cores - $1, 0),
		    used_ram_mb    = GREATEST(used_ram_mb - $2, 0)
		WHERE user_id = $3`,
		res.ReqCPU, res.ReqRAM, res.UserId,
	)
	return err
}




// // handleStop: running → stopped (manual oleh user)
// // Cancel timer, release quota, set cooldown next_start_at, publish stop.
// func (r *Repository) handleStop(ctx context.Context, tx *sql.Tx, workspaceId string, res *workspaceRow) error {
// 	// Cancel timer — selalu, tidak perlu guard status
// 	messagebroker.CancelAutoStop(workspaceId, r.redis.Inceptor)

// 	// Release quota
// 	if err := r.releaseQuota(ctx, tx, res); err != nil {
// 		return fmt.Errorf("gagal release quota: %w", err)
// 	}

// 	// var pausedAt, expiresAt time.Time
// 	// err := tx.QueryRowContext(ctx, `
// 	// 	SELECT  expires_at
// 	// 	FROM workspace_sessions
// 	// 	WHERE workspace_id = $1 AND status = 'running'
// 	// 	ORDER BY created_at DESC
// 	// 	LIMIT 1
// 	// `, workspaceId).Scan(&expiresAt)
// 	// if err != nil {
// 	// 	log.Printf("failed to resumed workspaces : %s",err.Error())
// 	// 	return fmt.Errorf("session paused tidak ditemukan: %w", err)
// 	// }


// 	// // Sisa waktu = waktu yang belum terpakai sebelum di-pause
// 	// remainingTime := expiresAt.Sub(pausedAt)
// 	// if remainingTime <= 0 {
// 	// 	// Waktu sudah habis saat di-pause, beri grace period minimal
// 	// 	remainingTime = 1 * time.Minute
// 	// }
// 	// Set cooldown next_start_at — ini yang membedakan manual stop vs auto stop
// 	// User harus tunggu cooldown sebelum bisa start lagi
// 	nextStartAt := time.Now().UTC().Add(cooldown)
// 	_, err := tx.ExecContext(ctx, `
// 		UPDATE workspace_sessions
// 		SET status       = 'stopped',
// 		    stopped_at   = NOW(),
// 		    paused_at    = NULL,
// 		    next_start_at = $1,
// 		    updated_at   = NOW()
// 		WHERE workspace_id = $2 AND status = 'running'
// 	`, nextStartAt, workspaceId)
// 	if err != nil {
// 		return fmt.Errorf("gagal update session: %w", err)
// 	}

// 	r.publishStop(ctx, workspaceId, res)

// 	return nil
// }
