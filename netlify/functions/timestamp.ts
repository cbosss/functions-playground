import { Handler } from "@netlify/functions";

export const handler: Handler = async (event, context) => {
  return {
    statusCode: 50,
    body: JSON.stringify({
      timestamp: Date.now(),
    }),
    headers: {
      "content-type": "application/json",
    },
  };
};
