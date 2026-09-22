import "@nagare-app/services";

async function getToken() {
  return localStorage.getItem("token");
}

async function setToken(token: string | null): Promise<void> {
  localStorage.setItem("token", token!);
}

window.bindings = {
  getToken,
  setToken,
};
