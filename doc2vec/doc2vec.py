import sys
from os import listdir, path

from pyknp import Juman
from gensim import models
from gensim.models.doc2vec import LabeledSentence
import logging
logging.basicConfig(format='%(asctime)s : %(levelname)s : %(message)s', level=logging.INFO)

class SentencesIterator():
    def __init__(self, generator_function):
        self.generator_function = generator_function
        self.generator = self.generator_function()

    def __iter__(self):
        # reset the generator
        self.generator = self.generator_function()
        return self

    def __next__(self):
        result = next(self.generator)
        if result is None:
            print("StopIteration")
            raise StopIteration
        else:

            return result


class Learner:

#    def _doc_to_sentence(self):
#        for d in self.docs:
#            try:
#                words = self._split_into_words(d[0])


#            except ValueError:
#                print("ValueError")
#                continue
#            print("Learning ... tagging: " + str(d[1]))
#            yield LabeledSentence(words=words, tags=d[0:-1])
    

    def learn(self, read_func):
        sentences = SentencesIterator(read_func)
        model = models.Doc2Vec(sentences, dm=0, vector_size=300, window=15, alpha=.025, min_alpha=.025, min_count=1, sample=1e-6, compute_loss=True)

        print('\n訓練開始')
        for epoch in range(20):
            print('Epoch: {}'.format(epoch + 1))
            print(model.iter)
            model.train(sentences, total_examples=model.corpus_count, epochs=model.iter)
            
            model.alpha -= (0.025 - 0.0001) / 19
            model.min_alpha = model.alpha
        
        print('\nsaving...')
        model.save('doc2vec_2.model')
        print('\nsaved')


    