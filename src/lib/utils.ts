import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"
<<<<<<< HEAD

=======
 
>>>>>>> parent of e08ff49... chore(client): remove next.js files
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
