"use client";
import { signIn } from "next-auth/react";

export default function Page() {
  async function loginUser(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    const formData = new FormData(e.currentTarget);

    try {
      const res = await signIn("credentials", {
        username: formData.get("username"),
        password: formData.get("password"),
        redirect: true,
        callbackUrl: "/dashboard",
      });
    } catch (error) {
      console.error(error);
    }
  }

  return (
    <div>
      <form onSubmit={loginUser}>
        <h1>Login</h1>
        <input type="text" name="username" />
        <input type="password" name="password" />
        <input type="submit" value="Login" />
      </form>
    </div>
  );
}
