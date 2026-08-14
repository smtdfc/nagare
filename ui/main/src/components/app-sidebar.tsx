'use client'

import * as React from 'react'

import { NavMain } from '@/components/nav-main'
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'
import {
  Settings2Icon,
  LifeBuoyIcon,
  SendIcon,
  FrameIcon,
  PieChartIcon,
  MapIcon,
  TerminalIcon,
  MessageCircle,
  WorkflowIcon,
  Plug2Icon,
} from 'lucide-react'

const data = {
  user: {
    name: 'shadcn',
    email: 'm@example.com',
    avatar: '/avatars/shadcn.jpg',
  },
  navMain: [
    {
      title: 'Chat',
      url: '/',
      icon: <MessageCircle />,
      items: [],
    },
    {
      title: 'Workflows',
      url: '#',
      icon: <WorkflowIcon />,
      items: [
        {
          title: 'Manage',
          url: '#',
        },
        {
          title: 'Logs',
          url: '#',
        },
      ],
    },
    {
      title: 'Plugins',
      url: '#',
      icon: <Plug2Icon />,
      items: [
        {
          title: 'Manage',
          url: '/plugin/overview',
        },
        {
          title: 'Logs',
          url: '#',
        },
      ],
    },
    {
      title: 'Settings',
      url: '#',
      icon: <Settings2Icon />,
      items: [
        {
          title: 'General',
          url: '/settings/general',
        },
        {
          title: 'LLM Provider',
          url: '/settings/llm-provider/overview',
        },
        {
          title: 'Limits',
          url: '#',
        },
      ],
    },
  ],
}
export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar variant="inset" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" render={<a href="#" />}>
              <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                <TerminalIcon className="size-4" />
              </div>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">Nagare Agent</span>
                <span className="truncate text-xs">v2 beta</span>
              </div>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
      </SidebarContent>
    </Sidebar>
  )
}
