import { useRef } from "react";

function defaultComparer<T>(prev: T, next: T, key: string) {
  return prev[key as keyof T] === next[key as keyof T];
}

export function useStableObject<T extends Record<keyof T, any>>(
  object: T,
  comparer: (prev: T, next: T, key: string) => unknown = defaultComparer,
): T {
  const ref = useRef(object);

  if (ref.current !== object) {
    const lastStableObjectKeys = Object.keys(ref.current);
    const nextObjectKeys = Object.keys(object);
    if (nextObjectKeys.length !== lastStableObjectKeys.length) {
      ref.current = object;
    } else {
      for (const key of lastStableObjectKeys) {
        if (!comparer(ref.current, object, key)) {
          ref.current = object;
          break;
        }
      }
    }
  }

  return ref.current;
}
