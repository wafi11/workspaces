#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"
WORKSPACE_ID="78d93aa4-c47e-4d30-9fcd-f7ef359e1703"

echo "=============================="
echo "  COLLAB FLOW TEST"
echo "=============================="

# ─── STEP 1: Login Owner ───────────────────────────────────────────
echo ""
echo "▶ Step 1 - Login sebagai Owner..."

OWNER_RES=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin86@gmail.com","password":"node1234"}')

OWNER_TOKEN=$(echo "$OWNER_RES" | grep -o '"data":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$OWNER_TOKEN" ]; then
  echo "❌ Gagal login owner. Response:"
  echo "$OWNER_RES"
  exit 1
fi

echo "✅ Owner token didapat: ${OWNER_TOKEN:0:40}..."

# ─── STEP 2: Login Collaborator ────────────────────────────────────
echo ""
echo "▶ Step 2 - Login sebagai Collaborator..."

COLLAB_RES=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin82@gmail.com","password":"node1234"}')

COLLAB_AUTH_TOKEN=$(echo "$COLLAB_RES" | grep -o '"token":"[^"]*"' | head -1 | cut -d'"' -f4)
COLLAB_USER_ID=$(echo "$COLLAB_RES" | grep -o '"user_id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$COLLAB_AUTH_TOKEN" ]; then
  echo "❌ Gagal login collaborator. Response:"
  echo "$COLLAB_RES"
  exit 1
fi

echo "✅ Collab token didapat: ${COLLAB_AUTH_TOKEN:0:40}..."
echo "✅ Collab user_id: $COLLAB_USER_ID"

# ─── STEP 3: Owner Invite Collaborator ─────────────────────────────
echo ""
echo "▶ Step 3 - Owner invite collaborator..."

INVITE_RES=$(curl -s -X POST "$BASE_URL/workspaces/$WORKSPACE_ID/invite" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OWNER_TOKEN" \
  -d "{\"user_id\":\"$COLLAB_USER_ID\",\"role\":\"editor\"}")

INVITE_TOKEN=$(echo "$INVITE_RES" | grep -o '"token":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$INVITE_TOKEN" ]; then
  echo "❌ Gagal invite collaborator. Response:"
  echo "$INVITE_RES"
  exit 1
fi

echo "✅ Invite token didapat: ${INVITE_TOKEN:0:40}..."

# ─── STEP 4: Collab tukar invite token → session token ─────────────
echo ""
echo "▶ Step 4 - Tukar invite token jadi session token..."

SESSION_RES=$(curl -s -X GET "$BASE_URL/workspaces/$WORKSPACE_ID/collab/token?token=$INVITE_TOKEN" \
  -H "Authorization: Bearer $COLLAB_AUTH_TOKEN" \
  -H "Content-Type: application/json")

SESSION_TOKEN=$(echo "$SESSION_RES" | grep -o '"token":"[^"]*"' | head -1 | cut -d'"' -f4)
WS_URL=$(echo "$SESSION_RES" | grep -o '"ws_url":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$SESSION_TOKEN" ]; then
  echo "❌ Gagal dapat session token. Response:"
  echo "$SESSION_RES"
  exit 1
fi

echo "✅ Session token didapat: ${SESSION_TOKEN:0:40}..."
echo "✅ WS URL: $WS_URL"

# ─── STEP 5: Validasi session token ────────────────────────────────
echo ""
echo "▶ Step 5 - Validasi session token..."

VALIDATE_RES=$(curl -s -X POST "$BASE_URL/workspaces/$WORKSPACE_ID/collab" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $COLLAB_AUTH_TOKEN" \
  -d "{\"token\":\"$SESSION_TOKEN\"}")

echo "Response: $VALIDATE_RES"

USER_ID=$(echo "$VALIDATE_RES" | grep -o '"userId":"[^"]*"' | head -1 | cut -d'"' -f4)
ROLE=$(echo "$VALIDATE_RES" | grep -o '"role":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$USER_ID" ]; then
  echo "❌ Validasi gagal."
  exit 1
fi

echo "✅ Validasi berhasil!"
echo "   user_id : $USER_ID"
echo "   role    : $ROLE"

# ─── SUMMARY ───────────────────────────────────────────────────────
echo ""
echo "=============================="
echo "  SUMMARY"
echo "=============================="
echo "Workspace ID  : $WORKSPACE_ID"
echo "Collab UserID : $COLLAB_USER_ID"
echo "Invite Token  : ${INVITE_TOKEN:0:40}..."
echo "Session Token : ${SESSION_TOKEN:0:40}..."
echo "WS URL        : $WS_URL"
echo ""
echo "🎉 Semua step berhasil!"