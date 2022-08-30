/** @jsx h */
import { Fragment, h } from "preact";
import { tw } from "@twind";
import { asset } from "$fresh/runtime.ts";
import Evaluator from "../islands/Evaluator.tsx";

export default function Home() {
  return (
    <Fragment>
      <head>
        <script src={asset("/wasm_exec.js")} />
        <script
          dangerouslySetInnerHTML={{
            __html: `const go = new Go();
      WebAssembly.instantiateStreaming(fetch("go.wasm"), go.importObject).then(
        (result) => {
          go.run(result.instance);
        }
      );`,
          }}
        />
        <style
          dangerouslySetInnerHTML={{
            __html: `
          html, body {height: 100%;}
          `,
          }}
        />
      </head>
      <div
        class={tw`
          bg-indigo-900 max-h-screen h-full flex flex-col py-4 items-center
        `}
      >
        <div
          class={tw`max-w-screen-xl px-4 space-y-8 h-full w-full flex flex-col`}
        >
          <h1 class={tw`text-4xl font-bold text-yellow-400`}>
            Gonkey Interpreter Playground!🐒
          </h1>
          <Evaluator />
        </div>
      </div>
    </Fragment>
  );
}
