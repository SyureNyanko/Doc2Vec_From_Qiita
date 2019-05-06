import glob
import json
import lxml.html
import os
from doc2vec import Learner
import re
import MeCab
from gensim.models.doc2vec import LabeledSentence

def ReadDocs():
    for path in sorted(glob.glob("pre2/*")):
        filename = os.path.basename(path)
        f = open(path, "r")
        line = f.readline()
        while line:
            yield LabeledSentence(words=line.split(' '), tags=filename.split('_')[0:-1])
            line = f.readline()
        f.close()


class PreDocsProcess: 
    def __init__(self):
        self.m = MeCab.Tagger ("-Owakati")


    def ReadLine(self, path):
        with open(path, encoding="utf-8") as f:
            for line in f:
                yield line

    def ExtractContents(self, path):
        for line in self.ReadLine(path):
            yield line

    def __iter__(self):
        f = False
        for path in sorted(glob.glob("dataset2/*"), key=os.path.getsize):
            i = 0
            print(path)
            for req_line in self.ExtractContents(path):
                json_dict = json.loads(req_line)
                i=i+1
                for p, j in enumerate(json_dict):
                    try:
                        raw_text = self.CleanContent(j["rendered_body"])
                        if raw_text == "" or raw_text is None:
                            continue
                        raw_text = self.Split(raw_text)
                        if raw_text is None:
                            continue
                        yield raw_text, path.split('/')[-1], i
                    except ValueError:
                        print("Error Value Error")
                        continue
                        

                
    def CleanContent(self, s):
        document = lxml.html.document_fromstring(s)
        raw_text = document.text_content()
#        raw_text=raw_text.replace(' ', '').replace('#', '').replace('\n', '').replace('@', '').replace('\u3000', '　')
        return raw_text

    def Split(self, text):
        result = self.m.parse(text)
        if result is None:
            return None
        return result.split(" ")



if __name__=='__main__':
    """
    ds = PreDocsProcess()
    index=0
    for d in ds:
        with open("./pre2/"+d[1]+"_" + str(d[2]), 'a') as f:
            l_in = [s for s in d[0] if not re.match('^[a-zA-Z0-9]$', s)]
            str_ = ' '.join(l_in)
            f.write(str_)
            f.write(' ')
        index = index + 1
    """
    l = Learner()
    l.learn(ReadDocs)