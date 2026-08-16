import type { ColumnDef } from "@tanstack/react-table";
import type { PluginInfo } from "@nagare-agent/dto";
import { Button } from "#/components/ui/button.tsx";
import { NutOff, Trash2 } from "lucide-react";

export const createColumns = (
  onDisable: (p: PluginInfo) => void,
  onDelete: (p: PluginInfo) => void,
): ColumnDef<PluginInfo>[] => {
  return [
    {
      accessorKey: "plugin_id",
      header: "Plugin ID",
      cell: ({ row }) => <span>{row.getValue("plugin_id")}</span>,
    },
    {
      accessorKey: "name",
      header: "Name",
      cell: ({ row }) => (
        <span className="font-semibold">{row.getValue("name")}</span>
      ),
    },
    {
      accessorKey: "version",
      header: "Version",
      cell: ({ row }) => {
        const val = row.getValue("version");
        return val as string;
      },
    },
    {
      accessorKey: "author",
      header: "Author",
      cell: ({ row }) => {
        const val = row.getValue("author");
        return (
          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-secondary text-secondary-foreground">
            {val as string}
          </span>
        );
      },
    },
    {
      id: "actions",
      header: "Actions",
      cell: ({ row }) => {
        const item = row.original;
        return (
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              className="h-8 px-2.5 text-xs flex items-center gap-1"
              onClick={() => onDisable(item)}
            >
              <NutOff className="size-3.5" />
              Disable
            </Button>
            <Button
              variant="destructive"
              size="sm"
              className="h-8 px-2.5 text-xs flex items-center gap-1"
              onClick={() => onDelete(item)}
            >
              <Trash2 className="size-3.5" />
              Delete
            </Button>
          </div>
        );
      },
    },
  ];
};
