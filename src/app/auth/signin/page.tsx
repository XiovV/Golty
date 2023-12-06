"use client";

import { redirect, useRouter } from "next/navigation";
import { signIn } from "next-auth/react";
import { useState } from "react";

export default function Page() {
  const router = useRouter();

  const [success, setSuccess] = useState(true);

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

      console.log("RES", res);
    } catch (e) {
      console.log(e);
    }
  }

  return (
    <div>
      {!success && <h1>Credentials are invalid</h1>}
      <form onSubmit={loginUser}>
        <h1>Login</h1>
        <input type="text" name="username" />
        <input type="password" name="password" />
        <input type="submit" value="Login" />
      </form>
    </div>
  );
}
