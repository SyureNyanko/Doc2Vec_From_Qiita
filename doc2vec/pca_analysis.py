import numpy as np
import pandas as pd
from gensim.models.doc2vec import Doc2Vec
from sklearn.decomposition import PCA
from MulticoreTSNE import MulticoreTSNE as TSNE
import matplotlib.pyplot as plt
import random

def get_taglist(model):
    taglist = []
    number = 0
    while True:
        tag = model.docvecs.index_to_doctag(number)
        if tag == number:
            break
        taglist.append(tag)
        number = number + 1
    return taglist


def plot_data(data, labels, filename):
    plt.figure(figsize=(12,9))
    # data = data[0:1000:5] #random.choices(data, k=10)
    plt.scatter(data[:,0], data[:,1], c=["w" for _ in labels])
    for d, l in zip(data,labels):
        plt.text(d[0], d[1], str(l), fontdict={"size":12, "color":'black'})
    plt.savefig(filename)

if __name__ == '__main__':
    model = Doc2Vec.load('doc2vec_2.model')
    tag_list = get_taglist(model)

    df_docvecs = pd.DataFrame()
    for tag in tag_list:
        df_docvecs[tag] = model.docvecs[tag]

    tsne_model = TSNE(n_jobs=4,
                  early_exaggeration=4,
                  n_components=2,
                  verbose=1,
                  random_state=2018,
                  n_iter=300)
    tsne_d2v = tsne_model.fit_transform(model.docvecs.vectors_docs)
    # tsne_d2v_df = pd.DataFrame(data=tsne_d2v, columns=["x", "y"])
    plot_data(tsne_d2v, tag_list, "test2.jpg")
