import { useState } from "react";
import { X } from "lucide-react";
import { Button } from "#/components/ui/button.tsx";
import { Input } from "#/components/ui/input.tsx";

type ModelListEditorProps = {
  models: string[];
  onChange: (models: string[]) => void;
};

export function ModelListEditor({ models, onChange }: ModelListEditorProps) {
  const [modelInput, setModelInput] = useState("");

  const handleAddModel = () => {
    const trimmed = modelInput.trim();
    if (!trimmed || models.includes(trimmed)) return;
    onChange([...models, trimmed]);
    setModelInput("");
  };

  const handleRemoveModel = (model: string) => {
    onChange(models.filter((m) => m !== model));
  };

  return (
    <div className="flex flex-col gap-1.5">
      <label className="text-sm font-medium">Models</label>
      <p className="text-xs text-muted-foreground">
        Add the model IDs this provider supports (chat and embedding).
      </p>
      <div className="flex gap-2">
        <Input
          placeholder="gpt-4o"
          value={modelInput}
          onChange={(e) => setModelInput(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              e.preventDefault();
              handleAddModel();
            }
          }}
        />
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={handleAddModel}
        >
          Add
        </Button>
      </div>
      {models.length > 0 && (
        <div className="mt-1.5 flex flex-wrap gap-1.5">
          {models.map((m) => (
            <span
              key={m}
              className="flex items-center gap-1 px-2 py-0.5 rounded-md bg-muted text-xs font-mono"
            >
              {m}
              <button
                type="button"
                onClick={() => handleRemoveModel(m)}
                className="text-muted-foreground hover:text-foreground"
                aria-label={`Remove model ${m}`}
              >
                <X className="w-3 h-3" />
              </button>
            </span>
          ))}
        </div>
      )}
    </div>
  );
}
