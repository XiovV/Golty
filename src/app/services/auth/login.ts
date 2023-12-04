import { Dispatch, SetStateAction } from "react";

interface TokenPair {
  accessToken: string;
  refreshToken: string;
}

export async function loginUser(e: React.FormEvent<HTMLFormElement>, setSuccess: Dispatch<SetStateAction<Boolean>>) {
  e.preventDefault();

  const formData = new FormData(e.currentTarget);

  const body = {
    username: formData.get("username"),
    password: formData.get("password"),
  };

  console.log(JSON.stringify(body));

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