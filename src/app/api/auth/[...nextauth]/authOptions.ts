import { AuthOptions } from "next-auth";
import CredentialsProvider from "next-auth/providers/credentials";

export const authOptions: AuthOptions = {
  callbacks: {
    async session({ session, token, user }) {
      console.log("session is called");

      console.log("SESSION", session)
      console.log("TOKEN", token)
      console.log("USER", user)
      return session
    }
  },
  providers: [CredentialsProvider({
    type: 'credentials',
    credentials: {},
    async authorize(credentials, req) {
      const { username, password } = credentials as {
        username: string;
        password: string;
      };

      const body = {
        username: username,
        password: password
      };

      try {
        const res = await fetch("http://localhost:8080/v1/users/login", {
          method: "POST",
          body: JSON.stringify(body),
          headers: new Headers({ "content-type": "application/json" }),
        });

        const token = await res.json();

        if (token.accessToken && token.refreshToken) {
          const user = {
            id: "this fixes the error",
            accessToken: token.accessToken,
            refreshToken: token.refreshToken
          };
          return user;
        } else {
          return null;
        }
      } catch (e) {
        console.log(e);
        return null;
      }
    }
  })],
  pages: {
    signIn: '/auth/signin'
  }
}