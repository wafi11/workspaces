import { Sun, SunMoon, Moon, Terminal } from "lucide-react"

export function relativeTime(iso: string): string {
  const diff = Math.floor((Date.now() - new Date(iso).getTime()) / 1000);
  if (diff < 60) return `${diff}s ago`;
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
  return `${Math.floor(diff / 86400)}d ago`;
}

export function formatDate(time: string) {
  if (!time) {
    return;
  }
  const date = Intl.DateTimeFormat("id-ID", {
    dateStyle: "full",
  });

  return date.format(new Date(time));
}





export function getSystemGreeting() {
  const hour = new Date().getHours()

  if (hour >= 5 && hour < 12) {
    return { 
      text: "System Online: Good Morning", 
      sub: "Everything is looking stable today.",
      icon: Sun 
    }
  }
  if (hour >= 12 && hour < 17) {
    return { 
      text: "High Traffic: Good Afternoon", 
      sub: "Don't forget to stay hydrated while coding.",
      icon: SunMoon 
    }
  }
  if (hour >= 17 && hour < 21) {
    return { 
      text: "Evening Deployment", 
      sub: "Winding down for the day?",
      icon: Moon 
    }
  }
  return { 
    text: "Night Mode: Happy Hacking", 
    sub: "The best code is written after midnight.",
    icon: Terminal 
  }
}