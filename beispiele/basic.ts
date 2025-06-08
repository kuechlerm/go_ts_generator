import { type } from "arktype";

export const Eins_Request_Schema = type({
  requiredString: "string > 0",
  optionalString: "string | undefined",
  requiredInt: "number > 0",
  optionalInt: "number | undefined",
  requiredBool: "true",
  optionalBool: "boolean | undefined",
});

export type Eins_Request = typeof Eins_Request_Schema.infer;

export const Eins_Response_Schema = type({
  responseString: "string > 0",
});

export type Eins_Response = typeof Eins_Response_Schema.infer;

export const Zwei_Request_Schema = type({
  optionalString: "string | undefined",
});

export type Zwei_Request = typeof Zwei_Request_Schema.infer;

export const Zwei_Response_Schema = type({
  responseString: "string > 0",
});

export type Zwei_Response = typeof Zwei_Response_Schema.infer;

export class RPC_Client {
  constructor(public base_url: string) {}

  async #do_fetch<TRequest, TResponse>(
    path: string,
    args: TRequest,
  ): Promise<{ result: TResponse | null; error: string | null }> {
    try {
      const result = await fetch(new URL(path, this.base_url).href, {
        method: "POST",
        body: JSON.stringify(args),
      });

      if (!result.ok) {
        console.error(
          `Fetch error: ${result.status} ${result.statusText} for ${path}`,
        );
        return {
          result: null,
          error: `Fetch error: ${result.status} ${result.statusText}`,
        };
      }

      const data = await result.json();

      return {
        result: data as TResponse,
        error: null,
      };
    } catch (error) {
      console.error(`Error during fetch for ${path}:`, error);

      return {
        result: null,
        error: error instanceof Error ? error.message : "Unknown error",
      };
    }
  }

  eins = (args: Eins_Request) =>
    this.#do_fetch<Eins_Request, Eins_Response>("/eins", args);
  zwei = (args: Zwei_Request) =>
    this.#do_fetch<Zwei_Request, Zwei_Response>("/zwei", args);
}
