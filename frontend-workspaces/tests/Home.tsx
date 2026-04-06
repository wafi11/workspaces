// // app/page.tsx
// import { PageContainer } from "@/components/layouts/PageContainer";
// import { Badge } from "@/components/ui/badge";
// import { Button } from "@/components/ui/button";
// import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
// import { Plus, LayoutGrid, GitBranch, Terminal, ScrollText, Server } from "lucide-react";
// import Link from "next/link";

// const workspaces = [
//   { id: 1, name: "api-gateway", lang: "GO", stack: "code-server · Go 1.22", status: "running", lastActive: "2h ago" },
//   { id: 2, name: "frontend-workspaces", lang: "TS", stack: "code-server · Node 20", status: "running", lastActive: "5h ago" },
//   { id: 3, name: "workspace-operator", lang: "GO", stack: "code-server · Go 1.22", status: "stopped", lastActive: "1d ago" },
// ];

// const activity = [
//   { label: "api-gateway started", time: "2h ago", type: "success" },
//   { label: "workspace-operator stopped", time: "1d ago", type: "warning" },
//   { label: "frontend-workspaces created", time: "2d ago", type: "success" },
//   { label: "Quota updated: 4 CPU, 8 GB", time: "3d ago", type: "default" },
// ];

// const statusVariant: Record<string, "default" | "secondary" | "destructive" | "outline"> = {
//   running: "default",
//   stopped: "secondary",
//   pending: "outline",
// };

// export function HomePage() {
//   return (
//     <PageContainer withSidebar={false}>
//       {/* Topbar */}
//       <div className="flex items-center justify-between mb-8 pb-5 border-b border-border">
//         <div className="flex flex-col gap-0.5">
//           <span className="text-[11px] text-muted-foreground uppercase tracking-widest font-mono">
//             Developer Platform
//           </span>
//           <h1 className="text-xl font-medium">Good morning, Wafi</h1>
//         </div>
//         <div className="flex gap-2">
//           <Button variant="outline" size="sm">
//             <Link href="/templates">
//               <LayoutGrid className="w-3.5 h-3.5" />
//               Templates
//             </Link>
//           </Button>
//           <Button size="sm">
//             <Link href="/workspaces/new">
//               <Plus className="w-3.5 h-3.5" />
//               New Workspace
//             </Link>
//           </Button>
//         </div>
//       </div>

//       {/* Stats */}
//       <div className="grid grid-cols-2 md:grid-cols-4 gap-2.5 mb-8">
//         {[
//           { label: "Total Workspaces", value: "3", sub: "2 running" },
//           { label: "CPU Usage", value: "1.2", sub: "cores / 4 allocated" },
//           { label: "Memory", value: "2.4 GB", sub: "of 8 GB limit" },
//           { label: "Storage", value: "18 GB", sub: "62% used" },
//         ].map((s) => (
//           <div key={s.label} className="bg-muted rounded-md p-4">
//             <p className="text-[11px] text-muted-foreground uppercase tracking-wider mb-1.5">{s.label}</p>
//             <p className="text-[22px] font-medium font-mono leading-none">{s.value}</p>
//             <p className="text-[11px] text-muted-foreground mt-1">{s.sub}</p>
//           </div>
//         ))}
//       </div>

//       {/* Main grid */}
//       <div className="grid grid-cols-1 lg:grid-cols-[1fr_300px] gap-6">
//         {/* Workspace list */}
//         <div>
//           <div className="flex items-center justify-between mb-3">
//             <span className="text-[11px] text-muted-foreground uppercase tracking-widest">Your Workspaces</span>
//             <Link href="/workspaces" className="text-xs text-muted-foreground hover:text-foreground transition-colors">
//               View all →
//             </Link>
//           </div>
//           <div className="flex flex-col gap-2">
//             {workspaces.map((ws) => (
//               <Link key={ws.id} href={`/workspaces/${ws.id}`}>
//                 <div className="flex items-center gap-3.5 p-4 bg-card border border-border rounded-xl hover:border-border/60 transition-colors cursor-pointer">
//                   <div className="w-9 h-9 rounded-md bg-muted border border-border flex items-center justify-center text-[12px] font-medium text-muted-foreground shrink-0 font-mono">
//                     {ws.lang}
//                   </div>
//                   <div className="flex-1 min-w-0">
//                     <p className="text-sm font-medium">{ws.name}</p>
//                     <p className="text-xs text-muted-foreground mt-0.5 truncate font-mono">
//                       ns/wafi-dev · {ws.stack}
//                     </p>
//                   </div>
//                   <div className="flex items-center gap-2.5 shrink-0">
//                     <Badge variant={statusVariant[ws.status]} className="text-[11px] capitalize">
//                       {ws.status}
//                     </Badge>
//                     <span className="text-[11px] text-muted-foreground">{ws.lastActive}</span>
//                   </div>
//                 </div>
//               </Link>
//             ))}
//             <Link href="/workspaces/new">
//               <div className="flex items-center justify-center p-4 border border-dashed border-border rounded-xl hover:border-border/60 transition-colors cursor-pointer opacity-60 hover:opacity-100">
//                 <span className="text-sm text-muted-foreground">+ Create new workspace</span>
//               </div>
//             </Link>
//           </div>
//         </div>

//         {/* Side panel */}
//         <div className="flex flex-col gap-4">
//           <Card>
//             <CardHeader className="pb-3">
//               <CardTitle className="text-[11px] text-muted-foreground uppercase tracking-widest font-normal">
//                 Quick Actions
//               </CardTitle>
//             </CardHeader>
//             <CardContent className="pt-0 flex flex-col gap-1">
//               {[
//                 { icon: Plus, label: "New from template", href: "/templates" },
//                 { icon: GitBranch, label: "Import from Git", href: "/workspaces/import" },
//                 { icon: Server, label: "Manage quota", href: "/settings/quota" },
//                 { icon: ScrollText, label: "View logs", href: "/logs" },
//               ].map(({ icon: Icon, label, href }) => (
//                 <Link key={label} href={href}>
//                   <button className="flex items-center gap-2.5 w-full px-2.5 py-2 rounded-md text-sm text-muted-foreground hover:bg-muted hover:text-foreground transition-colors text-left">
//                     <div className="w-7 h-7 rounded-md bg-muted border border-border flex items-center justify-center shrink-0">
//                       <Icon className="w-3.5 h-3.5" />
//                     </div>
//                     {label}
//                   </button>
//                 </Link>
//               ))}
//             </CardContent>
//           </Card>

//           <Card>
//             <CardHeader className="pb-3">
//               <CardTitle className="text-[11px] text-muted-foreground uppercase tracking-widest font-normal">
//                 Recent Activity
//               </CardTitle>
//             </CardHeader>
//             <CardContent className="pt-0">
//               <div className="flex flex-col">
//                 {activity.map((a, i) => (
//                   <div key={i} className="flex items-start gap-2.5 py-2 border-b border-border last:border-0">
//                     <div className={`w-1.5 h-1.5 rounded-full mt-1.5 shrink-0 ${
//                       a.type === "success" ? "bg-green-700 dark:bg-green-400" :
//                       a.type === "warning" ? "bg-amber-600 dark:bg-amber-400" :
//                       "bg-muted-foreground"
//                     }`} />
//                     <span className="text-xs text-muted-foreground flex-1">{a.label}</span>
//                     <span className="text-[11px] text-muted-foreground shrink-0">{a.time}</span>
//                   </div>
//                 ))}
//               </div>
//             </CardContent>
//           </Card>
//         </div>
//       </div>
//     </PageContainer>
//   );
// }