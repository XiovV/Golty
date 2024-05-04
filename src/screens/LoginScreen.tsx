import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import goltyLogo from "../assets/golty.svg";
import { Label } from "@/components/ui/label";

function LoginScreen() {
  return (
    <>
      <div className="min-h-screen flex flex-col items-center">
        <div className="flex gap-2 items-center mt-12">
          <img src={goltyLogo} alt="golty logo" width={80} height={80} />
          <p className="text-3xl font-semibold">Golty</p>
        </div>

        <div className="flex flex-col max-w-md w-full px-8 mt-12">
          <p className="font-semibold text-lg">Login to your account</p>
          <div className="grid max-w-md w-full items-center gap-1.5 mt-5">
            <Label htmlFor="username">Username</Label>
            <Input type="text" id="username" placeholder="Username" />
          </div>

          <div className="grid max-w-md w-full items-center gap-1.5 mt-5">
            <Label htmlFor="password">Password</Label>
            <Input type="password" id="password" placeholder="Password" />
          </div>

          <Button variant="secondary" className="w-full mt-7">
            <p className="font-medium">Login</p>
          </Button>
        </div>
      </div>
    </>
  );
}

export default LoginScreen;
