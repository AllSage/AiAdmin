import { AiAdmin } from 'langchain/llms/AiAdmin';
import * as readline from "readline";

async function main() {
  const AiAdmin = new AiAdmin({
    model: 'mistral'    
    // other parameters can be found at https://js.langchain.com/docs/api/llms_AiAdmin/classes/AiAdmin
  });

  const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
  });

  rl.question("What is your question: \n", async (user_input) => {
    const stream = await AiAdmin.stream(user_input);
  
    for await (const chunk of stream) {
      process.stdout.write(chunk);
    }
    rl.close();
  })
}

main();