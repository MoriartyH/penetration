import csv,openpyxl,time,threading
from googletrans import Translator
from openpyxl.styles import PatternFill


n = 0  # 漏洞计数器

host = []   #IP列表
port = []   #端口列表
name = []   #漏洞名称列表
risk = []   #风险等级列表
description = []    #漏洞说明列表
solution = []   #加固建议列表
cve = []    #CVE列表
columns= ['IP','端口','漏洞名称','风险等级','漏洞说明','加固建议','CVE']


def trans(text):        #谷歌翻译
    tran = Translator()
    word = tran.translate(text,dest='zh-CN').text
    return word


def riskmean(risks):        #风险等级翻译
    if risks == 'Critical':
        m = '紧急'
    elif risks == 'High':
        m = '高危'
    elif risks == 'Medium':
        m = '中危'
    elif risks == 'Low':
        m = '低危'

    risk.append(m)

def ex():                  #创建Excel表，并写入数据
    global n
    wv = openpyxl.Workbook()
    wv.save(filename='Nessus漏洞表.xlsx')
    #print("创建表，成功")
    wb = openpyxl.load_workbook('Nessus漏洞表.xlsx')
    sheet = wb.active
    for i in range(1, len(columns) + 1):  # 写入表头
        _ = sheet.cell(row=1, column=i, value=str(columns[i - 1]))

    for j in range(2,9):
        if j == 2: data = host
        if j == 3: data = port
        if j == 4: data = name
        if j == 5: data = risk
        if j == 6: data = description
        if j == 7: data = solution
        if j == 8: data = cve
        #print(data)
        for i in range(2, len(data) + 1):  # 写入数据
            _ = sheet.cell(row=i, column=j-1, value=str(data[i - 1]))
    sheet.title='扫描结果'
    wb.save('Nessus漏洞表.xlsx')

    Critical = PatternFill("solid", fgColor="FF0000")
    High = PatternFill("solid", fgColor="FFA500")
    Medium = PatternFill("solid", fgColor="FFFF00")
    Low = PatternFill("solid", fgColor="C0FF3E")
    for i in range(2,n+1):
        if sheet.cell(row=i,column=4).value == '紧急':
            sheet.cell(row=i, column=4).fill = Critical
        elif sheet.cell(row=i,column=4).value == '高危':
            sheet.cell(row=i, column=4).fill = High
        elif sheet.cell(row=i,column=4).value == '中危':
            sheet.cell(row=i, column=4).fill = Medium
        elif sheet.cell(row=i,column=4).value == '低危':
            sheet.cell(row=i, column=4).fill = Low
    wb.save('Nessus漏洞表.xlsx')

    print("保存成功,共" + str(n) + '个漏洞')
    print('程序开始' + time.strftime("%a %b %d %H:%M:%S %Y", time.localtime()))

def runs(i):
    global n
    if i[1] != 'CVE':
        if i[3] != 'None':
            host.append(i[4])
            port.append(i[6])
            name.append(i[7])
            riskmean(i[3])
            description.append(i[9])
            solution.append(i[10])
            cve.append(i[1])

            n += 1


if __name__ == '__main__':
    print('程序开始' + time.strftime("%a %b %d %H:%M:%S %Y", time.localtime()))
    with open("1.csv","r",encoding="utf-8") as csvfile:
        reader = csv.reader(csvfile)
        for i in reader:
            t1 = threading.Thread(target=runs, args=(i,))
            t1.start()
            #time.sleep(1)
            t1.join()
    ex()