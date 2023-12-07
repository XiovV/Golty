import NextAuth, { Session } from "next-auth"

declare module "next-auth" {
  interface Session {
    user: {
      accessToken: string
      refreshToken: string
    }
  }
}