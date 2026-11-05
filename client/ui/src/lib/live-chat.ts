export type LiveChatSnapshot = {
  sessionID: string;
  userText: string;
  assistantText: string;
  status: string;
  updatedAt: number;
};

const STORAGE_PREFIX = "nagare:live-chat:";
export const LIVE_CHAT_MAX_AGE_MS = 10 * 60 * 1000;

function getStorage(): Storage | null {
  if (typeof window === "undefined") return null;

  try {
    return window.sessionStorage;
  } catch {
    return null;
  }
}

function storageKey(sessionID: string) {
  return `${STORAGE_PREFIX}${sessionID}`;
}

export function saveLiveChatSnapshot(snapshot: LiveChatSnapshot) {
  const storage = getStorage();
  if (!storage) return;

  try {
    storage.setItem(storageKey(snapshot.sessionID), JSON.stringify(snapshot));
  } catch {
    // Storage can be full or disabled; losing the snapshot only means the
    // user loses the in-progress bubble after a reload, not the turn itself.
  }
}

export function loadLiveChatSnapshot(
  sessionID: string,
  now = Date.now(),
): LiveChatSnapshot | null {
  const storage = getStorage();
  if (!storage) return null;

  try {
    const raw = storage.getItem(storageKey(sessionID));
    if (!raw) return null;

    const snapshot = JSON.parse(raw) as LiveChatSnapshot;
    if (snapshot.sessionID !== sessionID || !snapshot.updatedAt) {
      clearLiveChatSnapshot(sessionID);
      return null;
    }

    if (now - snapshot.updatedAt > LIVE_CHAT_MAX_AGE_MS) {
      clearLiveChatSnapshot(sessionID);
      return null;
    }

    return snapshot;
  } catch {
    clearLiveChatSnapshot(sessionID);
    return null;
  }
}

export function clearLiveChatSnapshot(sessionID: string) {
  const storage = getStorage();
  if (!storage) return;

  try {
    storage.removeItem(storageKey(sessionID));
  } catch {
    // Ignore: a stale snapshot expires on its own after LIVE_CHAT_MAX_AGE_MS.
  }
}
