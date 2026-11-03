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
import { envConfig } from "@nagare-app/services";
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
  const [gateway, setGateway] = useState(envConfig.baseRestApiUrl);
  const [token, setToken] = useState("");
  const [showAdvanced, setShowAdvanced] = useState(false);

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
            <FieldLabel htmlFor="token">Token</FieldLabel>
            <Input
              id="token"
              type="text"
              placeholder="Your token ... "
              autoFocus
              value={token}
              onChange={(e) => {
                setToken(e.target.value);
              }}
              required
            />
            <FieldDescription>
              Gateway: <span className="font-mono">{gateway}</span>
              {" · "}
              <button
                type="button"
                className="underline underline-offset-2 hover:text-foreground"
                onClick={() => setShowAdvanced((visible) => !visible)}
              >
                {showAdvanced ? "Hide advanced" : "Advanced"}
              </button>
            </FieldDescription>
          </Field>
          {showAdvanced ? (
            <Field>
              <FieldLabel htmlFor="gateway">Gateway URL</FieldLabel>
              <Input
                id="gateway"
                type="text"
                placeholder="http://localhost:9832"
                value={gateway}
                onChange={(e) => {
                  setGateway(e.target.value);
                }}
              />
              <FieldDescription>
                Change this when the gateway runs on another host or port.
              </FieldDescription>
            </Field>
          ) : null}
          <Field>
            <Button type="submit">Continue</Button>
          </Field>
          <FieldDescription className="text-center">
            Don&apos;t have a token?{" "}
            <a href="/helps/how-to-get-a-token">Get a token</a>
          </FieldDescription>
        </FieldGroup>
      </form>
    </div>
  );
}
