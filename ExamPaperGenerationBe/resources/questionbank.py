import pandas as pd

# 读取文本文件内容
with open('question.txt', 'r', encoding='utf-8') as file:
    lines = file.readlines()

# 初始化题目列表
questions = []
current_question = ""
for line in lines:
    line = line.strip()
    if line:
        if line.startswith(('A.', 'B.', 'C.', 'D.')):
            current_question += " " + line
        else:
            if current_question:
                questions.append(current_question)
            current_question = line
if current_question:
    questions.append(current_question)

# 创建DataFrame
df = pd.DataFrame(questions, columns=['题目'])

# 保存为Excel文件
df.to_excel('questions.xlsx', index=False)