const KEY = 'rune_token'

export const auth = {
  get: (): string => localStorage.getItem(KEY) ?? '',
  set: (token: string) => localStorage.setItem(KEY, token),
  clear: () => localStorage.removeItem(KEY),
}
