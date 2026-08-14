import { Link, useMatches } from '@tanstack/react-router'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import React from 'react'

export function DynamicBreadcrumbs() {
  const matches = useMatches()

  return (
    <Breadcrumb>
      <BreadcrumbList>
        {matches.map((match, index) => {
          const breadcrumbTitle = (match.staticData as { breadcrumb?: string })
            ?.breadcrumb
          if (!breadcrumbTitle) return null

          const isLast = index === matches.length - 1

          return (
            <React.Fragment key={match.id}>
              {index > 0 && <BreadcrumbSeparator className="hidden md:block" />}
              <BreadcrumbItem className={isLast ? '' : 'hidden md:block'}>
                {isLast ? (
                  <BreadcrumbPage>{breadcrumbTitle}</BreadcrumbPage>
                ) : (
                  <BreadcrumbLink>
                    <Link to={match.pathname}>{breadcrumbTitle}</Link>
                  </BreadcrumbLink>
                )}
              </BreadcrumbItem>
            </React.Fragment>
          )
        })}
      </BreadcrumbList>
    </Breadcrumb>
  )
}
