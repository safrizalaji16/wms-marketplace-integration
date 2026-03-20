import type { ChangeEvent, FormEvent } from "react";
import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { ArrowRight, LockKeyhole, Mail } from "lucide-react";
import { LogoMark } from "../icons/LogoMark";
import { useAuthStore } from "../../store/useAuthStore";
import { login as loginRequest } from "../../services/auth";
import type { AuthState } from "../../store/useAuthStore";

export function LoginScreen() {
  const login = useAuthStore((state: AuthState) => state.login);
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const loginMutation = useMutation({
    mutationFn: loginRequest,
    onSuccess: (data) => {
      const userName = email.split("@")[0] || "operator";
      login({ userName, token: data.token });
    }
  });

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    loginMutation.mutate({ email, password });
  }

  function handleEmailChange(event: ChangeEvent<HTMLInputElement>) {
    setEmail(event.target.value);
  }

  function handlePasswordChange(event: ChangeEvent<HTMLInputElement>) {
    setPassword(event.target.value);
  }

  return (
    <main className="relative flex min-h-screen items-center justify-center overflow-hidden bg-halo p-4">
      <div className="absolute inset-x-0 top-0 h-80 bg-[radial-gradient(circle_at_top,rgba(62,99,244,0.22),transparent_50%)]" />
      <div className="w-full max-w-xl">
        <section className="panel p-6 sm:p-8">
          <div className="mx-auto max-w-md">
            <div className="inline-flex items-center gap-4 rounded-full bg-primary-50 px-4 py-3 text-primary-700">
              <LogoMark />
              <div>
                <p className="text-xl font-extrabold">WMSpaceIO</p>
              </div>
            </div>

            <div className="mt-8">
              <p className="text-xs font-black uppercase tracking-[0.3em] text-primary-500">Login Form</p>
              <h2 className="mt-3 text-3xl font-extrabold tracking-tight text-slate-900">
                Welcome back to WMSpaceIO
              </h2>
              <p className="mt-3 text-sm leading-6 text-slate-500">
                Use a warehouse operator account to access internal WMS actions.
              </p>
            </div>

            <form className="mt-8 space-y-4" onSubmit={handleSubmit}>
              <label className="block">
                <span className="mb-2 block text-sm font-semibold text-slate-600">Email</span>
                <div className="flex items-center gap-3 rounded-3xl border border-slate-200 bg-white px-4 py-3 shadow-card">
                  <Mail size={18} className="text-slate-400" />
                  <input
                    type="email"
                    value={email}
                    placeholder="Enter your email"
                    onChange={handleEmailChange}
                    className="w-full border-none bg-transparent text-sm text-slate-700 outline-none"
                  />
                </div>
              </label>

              <label className="block">
                <span className="mb-2 block text-sm font-semibold text-slate-600">Password</span>
                <div className="flex items-center gap-3 rounded-3xl border border-slate-200 bg-white px-4 py-3 shadow-card">
                  <LockKeyhole size={18} className="text-slate-400" />
                  <input
                    type="password"
                    value={password}
                    placeholder="Enter your password"
                    onChange={handlePasswordChange}
                    className="w-full border-none bg-transparent text-sm text-slate-700 outline-none"
                  />
                </div>
              </label>

              <button
                type="submit"
                disabled={loginMutation.isPending}
                className="flex w-full items-center justify-center gap-2 rounded-[24px] bg-primary-600 px-4 py-3.5 text-sm font-semibold text-white shadow-lg shadow-primary-500/25 transition hover:bg-primary-700"
              >
                {loginMutation.isPending ? "Signing In..." : "Sign In"}
                <ArrowRight size={18} />
              </button>
            </form>

            {loginMutation.isError ? (
              <div className="mt-4 rounded-[24px] border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-600">
                {loginMutation.error instanceof Error
                  ? loginMutation.error.message
                  : "Login failed"}
              </div>
            ) : null}
          </div>
        </section>
      </div>
    </main>
  );
}
