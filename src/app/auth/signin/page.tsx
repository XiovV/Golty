"use client";

import { redirect, useRouter } from "next/navigation";
import { useState } from "react";

interface TokenPair {
  accessToken: string;
  refreshToken: string;
}

export default function Page() {
  const router = useRouter();

  const [success, setSuccess] = useState(true);

  async function loginUser(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    const formData = new FormData(e.currentTarget);

    const body = {
      username: formData.get("username"),
      password: formData.get("password"),
    };

    try {
      const res = await fetch("http://localhost:8080/v1/users/login", {
        method: "POST",
        body: JSON.stringify(body),
        headers: new Headers({ "content-type": "application/json" }),
      });

      if (res.status == 401) {
        setSuccess(false);
        return;
      }

      const token: TokenPair = await res.json();

      localStorage.setItem("accessToken", token.accessToken);
      localStorage.setItem("refreshToken", token.refreshToken);

      router.push("/dashboard");
    } catch (error) {
      console.error(error);
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
