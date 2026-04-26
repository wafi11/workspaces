import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faHouse,
  faDesktop,
  faFileLines,
  faBell,
  faUsers,
  faGear,
  faKey,
  faUser,
  faSlidersH,
  faUserGroup,
} from "@fortawesome/free-solid-svg-icons";
import {
  faBell as faBellReg,
  faFileLines as faFileLinesReg,
  faUser as faUserReg,
} from "@fortawesome/free-regular-svg-icons";
import type { NavItem } from "@/types";

export const NAV: NavItem[] = [
  {
    label: "Home",
    to: "/home",
    icon: <FontAwesomeIcon icon={faHouse} />,
    roles: ["user", "admin"],
  },
  {
    label: "Workspaces",
    to: "/workspaces",
    icon: <FontAwesomeIcon icon={faDesktop} />,
    roles: ["user", "admin"],
  },
  {
    label: "Templates",
    to: "/templates",
    icon: <FontAwesomeIcon icon={faFileLinesReg} />,
    roles: ["user", "admin"],
  },
  {
    label: "Notifications",
    to: "/notifications",
    icon: <FontAwesomeIcon icon={faBellReg} />,
    roles: ["user", "admin"],
  },
  {
    label: "Users",
    to: "/users",
    icon: <FontAwesomeIcon icon={faUsers} />,
    roles: ["admin"],
    badge: "Admin",
  },
  {
    label: "Settings",
    to: "/settings",
    icon: <FontAwesomeIcon icon={faGear} />,
    roles: ["user", "admin"],
  },
];

export const NavSettings: NavItem[] = [
  {
    label: "PAT",
    to: "/settings/pat",
    icon: <FontAwesomeIcon icon={faKey} />,
    roles: ["user", "admin"],
  },
  {
    label: "Profile",
    to: "/settings/profile",
    icon: <FontAwesomeIcon icon={faUserReg} />,
    roles: ["user", "admin"],
  },
  {
    label: "Preferences",
    to: "/settings/preferences",
    icon: <FontAwesomeIcon icon={faSlidersH} />,
    roles: ["user", "admin"],
  },
  {
    label: "Members",
    to: "/settings/members",
    icon: <FontAwesomeIcon icon={faUserGroup} />,
    roles: ["admin"],
  },
];