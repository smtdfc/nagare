import * as React from "react";

import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ChatService, LLMProviderService } from "@nagare-app/services";
import type { LLMProvider, Session } from "@nagare-app/dtos";
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
  GalleryVerticalEndIcon,
  TrashIcon,
  MoreHorizontalIcon,
  EditIcon,
  ArrowRightFromLine,
  Search,
  BrushCleaning,
  Layers,
  Archive,
  Brain,
  Check,
  Loader2,
  SlidersHorizontal,
} from "lucide-react";

type ChatNavActionsProps = {
  session: Session;
  onSessionUpdated: (session: Session) => void;
  onSessionDeleted: () => void;
  onSessionDuplicated: (session: Session) => void;
};

function ChatLLMSelector({
  session,
  onSessionUpdated,
}: Pick<ChatNavActionsProps, "session" | "onSessionUpdated">) {
  const [isOpen, setIsOpen] = React.useState(false);
  const [providers, setProviders] = React.useState<LLMProvider[]>([]);
  const [models, setModels] = React.useState<string[]>([]);
  const [providerID, setProviderID] = React.useState(
    session.currentLLMProvider?.id ?? "",
  );
  const [model, setModel] = React.useState(session.currentLLMModel ?? "");
  const [isLoading, setIsLoading] = React.useState(false);
  const [isSaving, setIsSaving] = React.useState(false);

  React.useEffect(() => {
    setProviderID(session.currentLLMProvider?.id ?? "");
    setModel(session.currentLLMModel ?? "");
  }, [session]);

  React.useEffect(() => {
    let isMounted = true;
    void LLMProviderService.list()
      .then((items) => {
        if (isMounted) setProviders(items ?? []);
      })
      .catch(showToastError);
    return () => {
      isMounted = false;
    };
  }, []);

  const loadModels = async (id: string) => {
    setIsLoading(true);
    try {
      setModels((await LLMProviderService.getModels({ id })) ?? []);
    } catch (error) {
      showToastError(error);
    } finally {
      setIsLoading(false);
    }
  };

  React.useEffect(() => {
    if (providerID) {
      void loadModels(providerID);
    }
  }, [providerID]);

  const handleProviderChange = (value: string | null) => {
    if (!value) return;
    setProviderID(value);
    setModel("");
  };

  const handleSave = async () => {
    if (!providerID || !model) return;
    setIsSaving(true);
    try {
      const updatedSession = await ChatService.updateChatSessionLLMSettings(
        session.id,
        {
          currentLLMProvider: providerID,
          currentLLMModel: model,
        },
      );
      onSessionUpdated(updatedSession);
      setIsOpen(false);
    } catch (error) {
      showToastError(error);
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <Popover open={isOpen} onOpenChange={setIsOpen}>
      <PopoverTrigger
        render={
          <Button
            variant="ghost"
            size="sm"
            className="h-8 gap-1.5 px-2 data-open:bg-accent"
            aria-label="Change provider and model"
          />
        }
      >
        <SlidersHorizontal className="h-4 w-4" />
        <span className="hidden max-w-28 truncate sm:inline">
          {session.currentLLMProvider?.name ?? "AI"}
        </span>
      </PopoverTrigger>
      <PopoverContent className="w-80" align="end">
        <div className="flex flex-col gap-4">
          <div>
            <h3 className="font-medium">Chat model</h3>
            <p className="text-xs text-muted-foreground">
              Choose the provider and model for this session.
            </p>
          </div>
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium">Provider</label>
            <Select value={providerID} onValueChange={handleProviderChange}>
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Select a provider">
                  {providers.find((provider) => provider.id === providerID)
                    ?.name ?? session.currentLLMProvider?.name}
                </SelectValue>
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  {providers.map((provider) => (
                    <SelectItem key={provider.id} value={provider.id}>
                      {provider.name}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium">Model</label>
            <div className="relative">
              {isLoading && (
                <Loader2 className="absolute right-8 top-2.5 z-10 h-4 w-4 animate-spin text-muted-foreground" />
              )}
              <Select
                value={model}
                disabled={isLoading || !providerID || models.length === 0}
                onValueChange={(value) => setModel(value ?? "")}
              >
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Select a model" />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    {models.map((item) => (
                      <SelectItem key={item} value={item}>
                        {item}
                      </SelectItem>
                    ))}
                  </SelectGroup>
                </SelectContent>
              </Select>
            </div>
          </div>
          <Button
            onClick={handleSave}
            disabled={isSaving || !providerID || !model}
          >
            {isSaving ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Check className="h-4 w-4" />
            )}
            Save
          </Button>
        </div>
      </PopoverContent>
    </Popover>
  );
}

export function ChatNavActions({
  session,
  onSessionUpdated,
  onSessionDeleted,
  onSessionDuplicated,
}: ChatNavActionsProps) {
  const [isOpen, setIsOpen] = React.useState(false);
  const [isActionLoading, setIsActionLoading] = React.useState(false);

  const handleArchive = async () => {
    setIsActionLoading(true);
    try {
      const updatedSession = await ChatService.archiveChatSession(session.id, {
        isArchive: !session.isArchive,
      });
      onSessionUpdated(updatedSession);
      setIsOpen(false);
    } catch (error) {
      showToastError(error);
    } finally {
      setIsActionLoading(false);
    }
  };

  const handleDuplicate = async () => {
    setIsActionLoading(true);
    try {
      const duplicatedSession = await ChatService.duplicateChatSession(
        session.id,
      );
      onSessionDuplicated(duplicatedSession);
      setIsOpen(false);
    } catch (error) {
      showToastError(error);
    } finally {
      setIsActionLoading(false);
    }
  };

  const handleDelete = async () => {
    setIsActionLoading(true);
    try {
      await ChatService.deleteChatSession(session.id);
      onSessionDeleted();
      setIsOpen(false);
    } catch (error) {
      showToastError(error);
    } finally {
      setIsActionLoading(false);
    }
  };

  return (
    <div className="flex items-center gap-2 text-sm">
      <ChatLLMSelector session={session} onSessionUpdated={onSessionUpdated} />
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
                        onClick={() => void handleDuplicate()}
                      >
                        <CopyIcon /> <span>Duplicate</span>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                    <SidebarMenuItem>
                      <SidebarMenuButton
                        disabled={isActionLoading}
                        onClick={() => void handleDelete()}
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
                      <SidebarMenuButton
                        disabled={isActionLoading}
                        onClick={() => void handleArchive()}
                      >
                        <Archive />
                        <span>
                          {session.isArchive ? "Unarchive" : "Archive"}
                        </span>
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
