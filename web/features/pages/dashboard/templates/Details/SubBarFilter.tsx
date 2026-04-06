import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ChevronDown, SlidersHorizontal } from "lucide-react";

interface SubBarFilterProps {
  filters: readonly string[];
  selected: (x: string) => boolean;
  setSelected: (filter: string) => void;
}

export function SubBarFilterComponent({
  filters,
  selected,

  setSelected,
}: SubBarFilterProps) {
  return (
    <div className="flex items-center gap-4 mt-4 border-t py-4">
      {/* Desktop: tampil semua button */}
      <div className="hidden sm:flex items-center gap-2 flex-wrap">
        {filters.map((filter) => (
          <Button
            key={filter}
            variant={selected(filter) ? "default" : "outline"}
            size="sm"
            onClick={() => setSelected(filter)}
          >
            {filter}
          </Button>
        ))}
      </div>

      {/* Mobile: dropdown */}
      <div className="flex sm:hidden">
        <DropdownMenu>
          <DropdownMenuTrigger
            className={"gap-2 flex items-center bg-accent p-2 rounded-xl"}
          >
            <SlidersHorizontal className="w-3.5 h-3.5" />
            {selected(filters[0]) ? filters[0] : "Filter"}
            <ChevronDown className="w-3.5 h-3.5" />
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start">
            {filters.map((filter) => (
              <DropdownMenuItem
                key={filter}
                onClick={() => setSelected(filter)}
                className={selected(filter) ? "bg-accent" : ""}
              >
                {filter}
              </DropdownMenuItem>
            ))}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  );
}
