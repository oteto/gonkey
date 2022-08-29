// deno-lint-ignore-file
/** @jsx h */
import { ComponentChildren, h } from "preact";
import { useState } from "preact/hooks";
import { tw } from "@twind";

interface Window {
  tokenize: any;
  parse: any;
  eval: any;
}
declare var window: Window;

const initialCode = `let hello_world = "Hello World!🐒";
puts(hello_world);
`;

export default function Evaluator() {
  const [input, setInput] = useState(initialCode);
  const [result, setResult] = useState("");
  return (
    <div class={tw`flex items-center space-x-4 h-full`}>
      <div
        class={tw`flex-1 flex flex-col space-y-2 h-full`}
      >
        <label
          htmlFor="code"
          class={tw`text-lg text-yellow-400 font-bold`}
        >
          Code
        </label>
        <textarea
          id="code"
          class={tw`
            resize-none focus:outline-none focus:ring-2 focus:ring-yellow-400 
            border border-grey-500 rounded-lg bg-black text-white p-4 h-full
          `}
          onChange={(e) => {
            setInput(e.currentTarget.value);
          }}
        >
          {input}
        </textarea>
      </div>

      <div class={tw`flex flex-col space-y-4`}>
        <Button
          onClick={() => {
            setResult(window.tokenize(input));
          }}
        >
          {"→ Tokenize →"}
        </Button>
        <Button
          onClick={() => {
            setResult(window.parse(input));
          }}
        >
          {"→ Parse →"}
        </Button>
        <Button
          onClick={() => {
            setResult(window.eval(input));
          }}
        >
          {"→ Exec →"}
        </Button>
      </div>

      <div class={tw`flex-1 flex flex-col space-y-2 h-full`}>
        <label
          htmlFor="code"
          class={tw`text-lg text-yellow-400 font-bold`}
        >
          Result
        </label>
        <textarea
          readOnly={true}
          class={tw`
            resize-none focus:outline-none focus:ring-2 focus:ring-yellow-400 
            border border-grey-500 rounded-lg bg-black text-white p-4 h-full
          `}
        >
          {result}
        </textarea>
      </div>
    </div>
  );
}

type Props = {
  children: ComponentChildren;
  onClick: () => void;
};

const Button = ({ children, onClick }: Props) => {
  return (
    <button
      class={tw`
            py-2 px-4 text-md font-bold bg-white text-yellow-600
            border-2 border-yellow-400 rounded-md focus:outline-none
            transition duration-200
            hover:bg-yellow-400 hover:text-white 
            active:bg-yellow-500
          `}
      onClick={onClick}
    >
      {children}
    </button>
  );
};
