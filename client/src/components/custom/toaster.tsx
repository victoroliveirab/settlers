import { Toaster as Sonner, type ToasterProps } from "sonner";

export const Toaster = (props: ToasterProps) => {
  return (
    <Sonner position="top-right" style={{ position: "absolute", top: 0, right: 0 }} {...props} />
  );
};
