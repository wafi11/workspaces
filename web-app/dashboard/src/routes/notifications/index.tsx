import { useGetnotificationReceived, useGetnotificationRetreived } from '@/features/api/notifications'
import { MainContainer } from '@/features/layout/MainContainer'
import { TopbarAdmin } from '@/features/layout/TopbarDashboard'
import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { ChevronDown, Bell, Inbox } from 'lucide-react'
import { formatDate } from '@/utils/formatDate'
import { useAcceptOrDenied } from '@/features/api/workspace-collaboration'
import { useProfile } from '@/features/api'

export const Route = createFileRoute('/notifications/')({
  component: RouteComponent,
})

type Tab = 'inbox' | 'retrieved'

function RouteComponent() {
  const { data: dataReceived } = useGetnotificationReceived()
  const { data: dataRetreived } = useGetnotificationRetreived()
  const {data : profileData}  = useProfile()
  const [activeTab, setActiveTab] = useState<Tab>('inbox')
  const [expanded, setExpanded] = useState<string | null>(null)
  const {mutate} = useAcceptOrDenied()

  const list = activeTab === 'inbox' ? dataReceived : dataRetreived

  const toggle = (id: string) => setExpanded(prev => prev === id ? null : id)
  const handleAcc = ({data} : {data : {types : string,notification_id : string}})  => {
    mutate({data})
  }
  return (
    <MainContainer>
      <TopbarAdmin title="Notifications" />

      {/* Tabs */}
      <div className="flex gap-1 border-b"
        style={{ borderColor: 'var(--color-sidebar-border)' }}>
        {(['inbox', 'retrieved'] as Tab[]).map(tab => (
          <button
            key={tab}
            onClick={() => setActiveTab(tab)}
            className={`flex items-center gap-2 px-4 py-2 text-xs font-mono capitalize
                        transition-colors border-b-2 -mb-px
                        ${activeTab === tab
                          ? 'border-zinc-200 text-zinc-200'
                          : 'border-transparent text-zinc-500 hover:text-zinc-300'
                        }`}
          >
            {tab === 'inbox' ? <Bell size={12} /> : <Inbox size={12} />}
            {tab}
            {tab === 'inbox' && dataReceived?.length ? (
              <span className="bg-zinc-700 text-zinc-300 text-[10px] px-1.5 py-0.5 rounded-full">
                {dataReceived.length}
              </span>
            ) : null}
          </button>
        ))}
      </div>

      {/* List */}
      <div className="flex flex-col gap-2">
        {!list?.length && (
          <p className="text-xs font-mono text-zinc-600 py-8 text-center">
            No notifications
          </p>
        )}
        {list?.map(notif => (
          <div
            key={notif.id}
            className="overflow-hidden transition-colors"
            style={{
              background: 'var(--color-sidebar-bg)',
              border: `1px solid var(--color-sidebar-border)`,
            }}
          >
            {/* Header — clickable */}
            <button
              onClick={() => toggle(notif.id)}
              className="w-full flex items-center gap-3 px-4 py-3 text-left hover:bg-zinc-800/40 transition-colors"
            >
              {/* Unread dot */}
              <span className={`w-1.5 h-1.5 rounded-full shrink-0 ${notif.is_read ? 'bg-transparent' : 'bg-blue-400'}`} />

              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between gap-2">
                  <span className="text-xs font-mono text-zinc-200 truncate">{notif.title}</span>
                  <span className="text-[10px] font-mono text-zinc-600 shrink-0">
                    {formatDate(notif.created_at)}
                  </span>
                </div>
                <p className="text-[11px] font-mono text-zinc-500 truncate mt-0.5">
                  {notif.message}
                </p>
              </div>

              <ChevronDown
                size={14}
                className={`text-zinc-500 shrink-0 transition-transform ${expanded === notif.id ? 'rotate-180' : ''}`}
              />
            </button>

            {/* Expanded content */}
            {expanded === notif.id && (
              <div className="px-4 pb-4 pt-1 border-t"
                style={{ borderColor: 'var(--color-sidebar-border)' }}>

                {/* Metadata */}
                {notif.metadata && (
                  <div className="grid grid-cols-2 gap-2 mb-3">
                    {Object.entries(notif.metadata).map(([key, val]) => (
                      key !== 'invite_token' && (
                        <div key={key} className="flex flex-col gap-0.5">
                          <span className="text-[10px] font-mono text-zinc-600 uppercase tracking-widest">
                            {key.replace(/_/g, ' ')}
                          </span>
                          <span className="text-[11px] font-mono text-zinc-300 truncate">
                            {String(val)}
                          </span>
                        </div>
                      )
                    ))}
                  </div>
                )}

                {/* Actions — khusus INVITATION_COLLABORATOR */}
                {!notif.is_read && notif.user_id === profileData?.id && notif.notification_type === 'INVITATION_COLLABORATOR' && (
                  <div className="flex gap-2 mt-3">
                    <button onClick={() => handleAcc({
                      data : {
                        notification_id : notif.id,
                        types : "accept"
                      }
                    })} className="text-[11px] font-mono px-3 py-1.5 bg-zinc-200 text-zinc-900
                                       hover:bg-white rounded transition-colors">
                      Accept
                    </button>
                    <button onClick={() => handleAcc({
                      data : {
                        notification_id : notif.id,
                        types : "decline"
                      }
                    })} className="text-[11px] font-mono px-3 py-1.5 border border-zinc-700
                                       text-zinc-400 hover:text-zinc-200 rounded transition-colors">
                      Decline
                    </button>
                  </div>
                )}
              </div>
            )}
          </div>
        ))}
      </div>
    </MainContainer>
  )
}