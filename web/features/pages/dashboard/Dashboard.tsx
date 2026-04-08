"use client";

import { PageContainer } from "@/components/layouts";
import { TopBarAdmin } from "@/components/layouts/TopBarAdmin";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CardDetails } from "@/components/ui/card-details";
import {
  Activity,
  Users,
  HardDrive,
  ShieldAlert,
  Settings,
  Database,
  Cpu,
  ArrowUpRight,
  UserPlus,
} from "lucide-react";
import Link from "next/link";

// Data Dummy untuk Admin
const nodes = [
  {
    id: "n1",
    name: "node-master-01",
    status: "ready",
    cpu: "12%",
    ram: "4.2GB",
    ip: "192.168.1.10",
  },
  {
    id: "n2",
    name: "worker-intel-xeon",
    status: "ready",
    cpu: "65%",
    ram: "42GB",
    ip: "192.168.1.11",
  },
  {
    id: "n3",
    name: "worker-thinkpad-x250",
    status: "ready",
    cpu: "28%",
    ram: "6GB",
    ip: "192.168.1.12",
  },
  {
    id: "n4",
    name: "storage-nas",
    status: "maintenance",
    cpu: "2%",
    ram: "1GB",
    ip: "192.168.1.20",
  },
];

const systemAlerts = [
  {
    id: 1,
    msg: "High CPU usage on worker-intel-xeon",
    time: "5m ago",
    severity: "destructive",
  },
  {
    id: 2,
    msg: "New user registration: 'wafi_dev'",
    time: "1h ago",
    severity: "default",
  },
  {
    id: 3,
    msg: "Backup job 'daily-k8s' completed",
    time: "3h ago",
    severity: "secondary",
  },
];

export function AdminDashboardPage() {
  return (
    <>
      <TopBarAdmin
        description="System Overview: Cluster Alpha"
        title="Adminmode"
      >
        <Button variant="outline" size="sm">
          <Settings className="w-3.5 h-3.5 mr-2" />
          Cluster Config
        </Button>
        <Button size="sm">
          <UserPlus className="w-3.5 h-3.5 mr-2" />
          Invite User
        </Button>
      </TopBarAdmin>

      {/* Cluster Stats (Global) */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-2.5 mb-8">
        {[
          {
            label: "Total Users",
            value: "128",
            sub: "12 active now",
            icon: Users,
          },
          {
            label: "Cluster CPU",
            value: "32.4%",
            sub: "Avg across 4 nodes",
            icon: Cpu,
          },
          {
            label: "Cluster RAM",
            value: "54.2 GB",
            sub: "Used of 128 GB",
            icon: Activity,
          },
          {
            label: "Total PVC",
            value: "42",
            sub: "Storage provisioned",
            icon: Database,
          },
        ].map((s) => (
         <CardDetails key={s.label} icon={s.icon} label={s.label} sub={s.label} value={s.value}/>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-[1fr_320px] gap-6">
        {/* Node Management List */}
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <span className="text-[11px] text-muted-foreground uppercase tracking-widest font-semibold">
              Infrastructure Nodes
            </span>
            <span className="text-[10px] text-green-500 flex items-center gap-1 animate-pulse">
              ● Live Sync
            </span>
          </div>

          <div className="flex flex-col gap-2.5">
            {nodes.map((node) => (
              <div
                key={node.id}
                className="group flex items-center gap-4 p-4 bg-card border border-border rounded-xl hover:bg-muted/30 transition-all"
              >
                <div
                  className={`w-10 h-10 rounded-lg flex items-center justify-center shrink-0 border border-border ${node.status === "ready" ? "bg-primary/5" : "bg-destructive/5"}`}
                >
                  <HardDrive
                    className={`w-5 h-5 ${node.status === "ready" ? "text-primary" : "text-destructive"}`}
                  />
                </div>

                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <p className="text-sm font-semibold truncate">
                      {node.name}
                    </p>
                    <Badge
                      variant={
                        node.status === "ready" ? "default" : "destructive"
                      }
                      className="text-[9px] h-4 uppercase"
                    >
                      {node.status}
                    </Badge>
                  </div>
                  <p className="text-[11px] text-muted-foreground font-mono mt-0.5">
                    {node.ip}
                  </p>
                </div>

                {/* Resource Mini Indicators */}
                <div className="hidden md:flex items-center gap-6 text-right mr-4">
                  <div className="space-y-1">
                    <p className="text-[9px] text-muted-foreground uppercase">
                      CPU
                    </p>
                    <p className="text-xs font-mono">{node.cpu}</p>
                  </div>
                  <div className="space-y-1">
                    <p className="text-[9px] text-muted-foreground uppercase">
                      RAM
                    </p>
                    <p className="text-xs font-mono">{node.ram}</p>
                  </div>
                </div>

                <Button
                  variant="ghost"
                  size="icon"
                  className="shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
                >
                  <ArrowUpRight className="w-4 h-4" />
                </Button>
              </div>
            ))}
          </div>
        </div>

        {/* Right Sidebar: Alerts & Admin Actions */}
        <div className="space-y-4">
          {/* System Alerts */}
          <Card className="border-destructive/20 bg-destructive/5">
            <CardHeader className="pb-3 flex-row items-center justify-between space-y-0">
              <CardTitle className="text-[11px] text-destructive uppercase tracking-widest font-bold flex items-center gap-2">
                <ShieldAlert className="w-3.5 h-3.5" />
                System Alerts
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-0 space-y-3">
              {systemAlerts.map((alert) => (
                <div
                  key={alert.id}
                  className="flex flex-col gap-1 border-b border-destructive/10 last:border-0 pb-2 last:pb-0"
                >
                  <p className="text-xs font-medium text-foreground leading-tight">
                    {alert.msg}
                  </p>
                  <span className="text-[10px] text-muted-foreground">
                    {alert.time}
                  </span>
                </div>
              ))}
            </CardContent>
          </Card>

          {/* Admin Tools */}
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-[11px] text-muted-foreground uppercase tracking-widest font-normal">
                Administrative Tools
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-0 flex flex-col gap-1">
              {[
                { label: "Namespace Limits", href: "/admin/limits" },
                { label: "Billing & Invoices", href: "/admin/billing" },
                { label: "Audit Logs", href: "/admin/audit" },
                { label: "Global Templates", href: "/admin/templates" },
              ].map((item) => (
                <Link key={item.label} href={item.href} className="w-full">
                  <button className="w-full text-left px-3 py-2 rounded-md text-xs text-muted-foreground hover:bg-muted hover:text-foreground transition-colors border border-transparent hover:border-border">
                    {item.label}
                  </button>
                </Link>
              ))}
            </CardContent>
          </Card>
        </div>
      </div>
    </>
  );
}
