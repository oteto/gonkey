import { MiddlewareHandlerContext } from "$fresh/server.ts";

export async function handler(
  req: Request,
  ctx: MiddlewareHandlerContext,
) {
  const url = new URL(req.url);
  if (url.pathname === "/") {
    console.log(ctx.localAddr);
  }
  const resp = await ctx.next();
  return resp;
}
