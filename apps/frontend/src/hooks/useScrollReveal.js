import { useEffect, useRef } from 'react'

/**
 * Attach to any element — adds 'visible' class when it enters the viewport.
 * @param {number} threshold – how much of the element must be visible (0–1)
 * @param {number} delay – extra delay in ms before adding the class
 */
export function useScrollReveal(threshold = 0.1, delay = 0) {
  const ref = useRef(null)

  useEffect(() => {
    const el = ref.current
    if (!el) return

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setTimeout(() => el.classList.add('visible'), delay)
          observer.unobserve(el)
        }
      },
      { threshold }
    )

    observer.observe(el)
    return () => observer.disconnect()
  }, [threshold, delay])

  return ref
}
