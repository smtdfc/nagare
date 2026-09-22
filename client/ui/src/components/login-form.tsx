"use client";

import { cn } from "cn";

import { Button } from "@/components/ui/button";
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { GalleryVerticalEndIcon } from "lucide-react";
import { type SyntheticEvent, useState } from "react";

export type LoginInput = {
  gateway?: string;
  token?: string;
};
type LoginFormProps = {
  onSubmit(i: LoginInput): void;
};

export function LoginForm({ onSubmit }: LoginFormProps) {
  const [gateway, setGateway] = useState<string>();
  const [token, setToken] = useState<string>();

  const handleSubmit = (event: SyntheticEvent) => {
    event.preventDefault();
    onSubmit({
      gateway,
      token,
    });
  };

  return (
    <div className={cn("flex flex-col gap-6")}>
      <form onSubmit={handleSubmit}>
        <FieldGroup>
          <div className="flex flex-col items-center gap-2 text-center">
            <a
              href="#"
              className="flex flex-col items-center gap-2 font-medium"
            >
              <div className="flex size-8 items-center justify-center rounded-md">
                <GalleryVerticalEndIcon className="size-6" />
              </div>
              <span className="sr-only">Acme.</span>
            </a>
            <h1 className="text-xl font-bold">Welcome to Nagare.</h1>
          </div>
          <Field>
            <FieldLabel htmlFor="gateway">Gateway URL</FieldLabel>
            <Input
              id="gateway"
              type="text"
              placeholder="Your gateway ... "
              defaultValue={gateway}
              onChange={(e) => {
                setGateway(e.target.value as string);
              }}
              required
            />
          </Field>
          <Field>
            <FieldLabel htmlFor="token">Token</FieldLabel>
            <Input
              id="token"
              type="text"
              placeholder="Your token ... "
              defaultValue={token}
              onChange={(e) => {
                setToken(e.target.value as string);
              }}
              required
            />
          </Field>
          <Field>
            <Button type="submit">Continue</Button>
          </Field>
          <FieldDescription className="text-center">
            Don&apos;t have a token? <a href="#">Get a token</a>
          </FieldDescription>
        </FieldGroup>
      </form>
    </div>
  );
}
