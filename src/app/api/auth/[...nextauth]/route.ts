import NextAuth, { NextAuthOptions, User } from "next-auth";
import CredentialsProvider from "next-auth/providers/credentials";

const authConfig = NextAuth({
  callbacks: {
    session({ session, token, user }) {
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
      }


      try {

        const res = await fetch("http://localhost:8080/v1/users/login", {
          method: "POST",
          body: JSON.stringify(body),
          headers: new Headers({ "content-type": "application/json" }),
        });

        const token = await res.json();

        const user = {
          accessToken: token.accessToken,
          refreshToken: token.refreshToken
        }

        return user;

      } catch (e) {
        console.log(e)
      }
    }

  })],
  pages: {
    signIn: '/auth/signin'
  }
})

export { authConfig as GET, authConfig as POST, authConfig as authConfig }