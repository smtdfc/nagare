import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Textarea } from "./ui/textarea";
import { useEffect, useState } from "react";
import { GetHost } from "@nagare-agent/service-bindings";

type AuthFormProps = {
  onSubmit?: (e: { host: string; token: string }) => void;
  onError?: (e: React.FormEvent) => void;
};

export function AuthForm({ onSubmit, onError, ...props }: AuthFormProps) {
  const [host, setHost] = useState("");
  const [token, setToken] = useState("");

  useEffect(() => {
    const getHost = async () => {
      const d = await GetHost();
      setHost(d);
    };
    getHost();
  }, []);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (onSubmit) onSubmit({ host, token });
  };

  return (
    <div className={cn("flex flex-col gap-6")} {...props}>
      <form onSubmit={handleSubmit}>
        <FieldGroup>
          <div className="flex flex-col items-center gap-2 text-center">
            <h1 className="text-xl font-bold">Welcome to Nagare Agent.</h1>
          </div>
          <Field>
            <FieldLabel htmlFor="host">Host</FieldLabel>
            <Input
              id="host"
              type="text"
              defaultValue={host}
              onChange={(e) => setHost(e.target.value)}
              placeholder=""
              disabled
            />
          </Field>
          <Field>
            <FieldLabel htmlFor="email">Token</FieldLabel>
            <Textarea
              id="token"
              rows={30}
              placeholder="Your token"
              defaultValue={token}
              onChange={(e) => setToken(e.target.value)}
              required
            />
          </Field>
          <Field>
            <Button type="submit">Continue</Button>
          </Field>
        </FieldGroup>
      </form>
      <FieldDescription className="px-6 text-center">
        You don't have a token?
      </FieldDescription>
    </div>
  );
}
