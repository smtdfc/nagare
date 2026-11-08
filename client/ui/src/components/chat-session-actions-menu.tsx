import * as React from "react";

import { Button } from "@/components/ui/button";
import type { Session } from "@nagare-app/dtos";
import { showToastError } from "#/lib/toast";
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
} from "@/components/ui/sidebar";
import {
  LinkIcon,
  CopyIcon,
  TrashIcon,
  MoreHorizontalIcon,
  EditIcon,
  Search,
  Layers,
  Archive,
  Brain,
} from "lucide-react";

type ChatSessionActionsMenuProps = {
  session: Session;
  onArchive: () => Promise<void>;
  onDuplicate: () => Promise<void>;
  onDelete: () => Promise<void>;
  onCompact: () => Promise<void>;
};

export function ChatSessionActionsMenu({
  session,
  onArchive,
  onDuplicate,
  onDelete,
  onCompact,
}: ChatSessionActionsMenuProps) {
  const [isOpen, setIsOpen] = React.useState(false);
  const [isActionLoading, setIsActionLoading] = React.useState(false);

  const runAction = async (action: () => Promise<void>) => {
    setIsActionLoading(true);
    try {
      await action();
      setIsOpen(false);
    } catch (error) {
      showToastError(error);
    } finally {
      setIsActionLoading(false);
    }
  };

  return (
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
                    <SidebarMenuButton
                      disabled={isActionLoading}
                      onClick={() => void runAction(onDuplicate)}
                    >
                      <CopyIcon /> <span>Duplicate</span>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                  <SidebarMenuItem>
                    <SidebarMenuButton
                      disabled={isActionLoading}
                      onClick={() => void runAction(onDelete)}
                    >
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
                      <Search /> <span>Search in chat</span>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                  <SidebarMenuItem>
                    <SidebarMenuButton
                      disabled={isActionLoading}
                      onClick={() => void runAction(onCompact)}
                    >
                      <Layers /> <span>Compact</span>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                  <SidebarMenuItem>
                    <SidebarMenuButton
                      disabled={isActionLoading}
                      onClick={() => void runAction(onArchive)}
                    >
                      <Archive />
                      <span>{session.isArchive ? "Unarchive" : "Archive"}</span>
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
  );
}
