from langchain.llms import AiAdmin

input = input("What is your question?")
llm = AiAdmin(model="llama3")
res = llm.predict(input)
print (res)
