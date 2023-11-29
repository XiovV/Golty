import NextAuth from "next-auth";

import CredentialsProvider from "next-auth/providers/credentials"

declare module "next-auth" {
  interface Session {
    user: {
      id: string;
      username: string;

    }
  }

  interface User {
    id: string;
    username: string;
  }
}

const authConfig = NextAuth({
  session: {
    strategy: "jwt"
  },
  providers: [CredentialsProvider({
    type: 'credentials',
    credentials: {},
    authorize(credentials, req) {
      const { username, password } = credentials as {
        username: string;
        password: string;
      };

      if (username !== "admin" || password !== "admin") {
        return null
      }

      return { id: "1", username: "Admin" }
    }
  })],
  pages: {
    signIn: '/auth/signin'
  }
})

export { authConfig as GET, authConfig as POST }