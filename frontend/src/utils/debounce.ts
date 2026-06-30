/**
 * Creates a debounced version of a function that delays invocation
 * until after `delay` ms have elapsed since the last call.
 *
 * @param fn  - The function to debounce
 * @param delay - Delay in milliseconds
 * @returns A debounced function with a .cancel() method
 */
export function debounce<T extends (...args: any[]) => any>(
  fn: T,
  delay: number,
): ((...args: Parameters<T>) => void) & { cancel: () => void } {
  let timer: ReturnType<typeof setTimeout> | null = null

  const debounced = (...args: Parameters<T>) => {
    if (timer !== null) {
      clearTimeout(timer)
    }
    timer = setTimeout(() => {
      timer = null
      fn(...args)
    }, delay)
  }

  debounced.cancel = () => {
    if (timer !== null) {
      clearTimeout(timer)
      timer = null
    }
  }

  return debounced
}
