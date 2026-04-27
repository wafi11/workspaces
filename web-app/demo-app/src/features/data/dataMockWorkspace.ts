import type { Workspace } from "@/types";

export const mockWorkspaces: Workspace[] = [

  {

    id: "1",

    name: "Workspace 1",

    template: "Python Data Science",

    description: "A Python data science environment with Jupyter, pandas, and scikit-learn.",

    icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/python/python-original.svg",

    url: "/workspace/1",

    status: "running",

    uptime: "2h 30m"

  },

  {

    id: "2",

    name: "Workspace 2",

    template: "Node.js Express",

    description: "A Node.js environment with Express, MongoDB, and Mongoose.",

    icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nodejs/nodejs-original.svg",

    url: "/workspace/2",

    status: "stopped",

    uptime: "0h 45m"

  },

  {

    id: "3",

    name: "Workspace 3",

    template: "React Frontend",

    description: "A React environment with TypeScript and Redux.",

    icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/react/react-original.svg",

    url: "/workspace/3",

    status: "running",

    uptime: "5h 10m"

  },

  {

    id: "4",

    name: "Workspace 4",

    template: "Go Backend",

    description: "A Go environment with Gin and GORM.",

    icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original.svg",

    url: "/workspace/4",

    status: "stopped",

    uptime: "1h 20m"

  },

  {

    id: "5",

    name: "Workspace 5",

    template: "Java Spring",

    description: "A Java Spring environment with Spring Boot and Hibernate.",

    icon: "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/java/java-original.svg",

    url: "/workspace/5",

    status: "running",

    uptime: "3h 5m"

  }

]


