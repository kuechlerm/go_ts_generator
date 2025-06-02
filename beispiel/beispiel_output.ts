// generated code - do not edit
import { type } from "arktype";

export const BeispielAnlegen_Request_Schema = type({
  name: "3 <= string <= 100",
});
export type BeispielAnlegen_Request =
  typeof BeispielAnlegen_Request_Schema.infer;

export const BeispielAnlegen_Response_Schema = type({
  id: "number",
});
export type BeispielAnlegen_Response =
  typeof BeispielAnlegen_Response_Schema.infer;

export const BeispielAendern_Request_Schema = type({
  name: "3 <= string <= 100",
});
export type BeispielAendern_Request =
  typeof BeispielAendern_Request_Schema.infer;

export const BeispielAendern_Response_Schema = type({
  id: "number",
});
export type BeispielAendern_Response =
  typeof BeispielAendern_Response_Schema.infer;

export class RPC_Client {
  constructor(public base_url: string) {}

  async #do_fetch<T>(path: string, args: T) {
    try {
      const result = await fetch(new URL(path, this.base_url).href, {
        method: "POST",
        body: JSON.stringify(args),
      });

      if (!result.ok) {
        throw new Error("fetch result not ok");
      }

      try {
        const data = await result.json();

        return data as T;
      } catch (parse_error) {
        throw parse_error;
      }
    } catch (fetch_error) {
      throw fetch_error;
    }
  }

  async beispiel_anlegen(args: BeispielAnlegen_Request) {
    // todo: validieren?
    return await this.#do_fetch<BeispielAnlegen_Response>(
      "/beispielanlegen",
      args,
    );
  }

  async beispiel_aendern(args: BeispielAendern_Request) {
    // todo: validieren?
    return await this.#do_fetch<BeispielAendern_Response>(
      "/beispielaendern",
      args,
    );
  }
}
