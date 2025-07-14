import { Api } from "./api";

async function getTokenFromGoogle(access_token: string) {
  return Api.post("/oauth/google/token", { access_token })
}

const oauthService = {
  getTokenFromGoogle
}

export { oauthService }
