import * as React from "react";

import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarSeparator,
} from "@/components/ui/sidebar";
import {
  Settings2Icon,
  LinkIcon,
  CopyIcon,
  CornerUpRightIcon,
  Trash2Icon,
  CornerUpLeftIcon,
  ChartLineIcon,
  GalleryVerticalEndIcon,
  TrashIcon,
  BellIcon,
  ArrowUpIcon,
  ArrowDownIcon,
  MoreHorizontalIcon,
  EditIcon,
  DotIcon,
  ArrowLeftFromLine,
  ArrowRightFromLine,
  Search,
  BrushCleaning,
  Layers,
  Archive,
  Brain,
} from "lucide-react";

const data = [
  [
    {
      label: "Rename",
      icon: <EditIcon />,
    },
    {
      label: "Copy Link",
      icon: <LinkIcon />,
    },
    {
      label: "Duplicate",
      icon: <CopyIcon />,
    },
    {
      label: "Delete",
      icon: <Trash2Icon />,
    },
  ],
  [
    {
      label: "Undo",
      icon: <CornerUpLeftIcon />,
    },
    {
      label: "Version History",
      icon: <GalleryVerticalEndIcon />,
    },
    {
      label: "Notifications",
      icon: <BellIcon />,
    },
  ],
  [
    {
      label: "Import",
      icon: <ArrowUpIcon />,
    },
    {
      label: "Export",
      icon: <ArrowDownIcon />,
    },
  ],
];
export function ChatNavActions() {
  const [isOpen, setIsOpen] = React.useState(false);

  return (
    <div className="flex items-center gap-2 text-sm">
      <Popover open={isOpen} onOpenChange={setIsOpen}>
        <PopoverTrigger
          render={
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 data-open:bg-accent"
            />
          }
        >
          <MoreHorizontalIcon />
        </PopoverTrigger>
        <PopoverContent
          className="w-56 overflow-hidden rounded-lg p-0"
          align="end"
        >
          <Sidebar collapsible="none" className="bg-transparent">
            <SidebarContent>
              <SidebarGroup className="border-b last:border-none">
                <SidebarGroupContent className="gap-0">
                  <SidebarMenu>
                    <SidebarMenuItem>
                      <SidebarMenuButton>
                        <EditIcon /> <span>Rename</span>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                    <SidebarMenuItem>
                      <SidebarMenuButton>
                        <LinkIcon /> <span>Copy link</span>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                    <SidebarMenuItem>
                      <SidebarMenuButton>
                        <CopyIcon /> <span>Duplicate</span>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                    <SidebarMenuItem>
                      <SidebarMenuButton>
                        <TrashIcon /> <span>Delete</span>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  </SidebarMenu>
                </SidebarGroupContent>
              </SidebarGroup>
              <SidebarGroup className="border-b last:border-none">
                <SidebarGroupContent className="gap-0">
                  <SidebarMenu>
                    <SidebarMenuItem>
                      <SidebarMenuButton>
                        <GalleryVerticalEndIcon /> <span>Checkpoints</span>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                    <SidebarMenuItem>
                      <SidebarMenuButton>
                        <ArrowRightFromLine /> <span>Restore to point</span>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                    <SidebarMenuItem>
                      <SidebarMenuButton>
                        <Search /> <span>Search in chat</span>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                    <SidebarMenuItem>
                      <SidebarMenuButton>
                        <BrushCleaning /> <span>Reset this chat</span>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                    <SidebarMenuItem>
                      <SidebarMenuButton>
                        <Layers /> <span>Compact</span>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                    <SidebarMenuItem>
                      <SidebarMenuButton>
                        <Archive /> <span>Archive</span>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                    <SidebarMenuItem>
                      <SidebarMenuButton>
                        <Brain /> <span>Add to knowledge pool</span>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  </SidebarMenu>
                </SidebarGroupContent>
              </SidebarGroup>
            </SidebarContent>
          </Sidebar>
        </PopoverContent>
      </Popover>
    </div>
  );
}
