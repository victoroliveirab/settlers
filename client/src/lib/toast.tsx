import { Check, Info, XCircle } from "lucide-react";
import { toast as sonnerToast } from "sonner";

import { cn } from "./utils";

export interface IToastProps {
  id: string | number;
  type?: "success" | "error" | "info";
  title: string;
  description: string;
  button?: {
    label: string;
    onClick: () => void;
  };
}

const Toast = (props: IToastProps) => {
  const { title, description, button, id, type } = props;

  let icon: React.ReactElement | undefined;
  if (type === "info") {
    icon = <Info size={16} />;
  } else if (type === "success") {
    icon = <Check size={16} />;
  } else if (type === "error") {
    icon = <XCircle size={16} />;
  }

  return (
    <div
      className={cn(
        "flex rounded-lg bg-white shadow-lg ring-1 ring-black/5 w-80 items-center p-4",
        {
          "bg-red-100": type === "error",
        },
      )}
    >
      <div className="flex flex-1 items-center">
        <div className="w-full">
          <p
            className={cn("text-sm font-medium text-gray-950 flex items-center gap-1", {
              "text-red-600": type === "error",
            })}
          >
            {icon} <span>{title}</span>
          </p>
          <p
            className={cn("mt-1 text-sm text-gray-500", {
              "text-gray-900": type === "error",
            })}
          >
            {description}
          </p>
        </div>
      </div>
      {!!button && (
        <div className="ml-5 shrink-0 rounded-md text-sm font-medium text-indigo-600 hover:text-indigo-500 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:outline-hidden">
          <button
            className="rounded bg-indigo-50 px-3 py-1 text-sm font-semibold text-indigo-600 hover:bg-indigo-100"
            onClick={() => {
              button.onClick();
              sonnerToast.dismiss(id);
            }}
          >
            {button.label}
          </button>
        </div>
      )}
    </div>
  );
};

export function toast(toast: Omit<IToastProps, "id">) {
  return sonnerToast.custom((id) => (
    <Toast
      id={id}
      title={toast.title}
      description={toast.description}
      button={toast.button}
      type={toast.type}
    />
  ));
}
