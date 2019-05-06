from qiita_doc2vec import Docs
import unittest

class LearnerTest(unittest.TestCase):
    def setUp(self):
        self.docs = Docs()
        # 初期化処理
        pass

    def tearDown(self):
        # 終了処理
        pass

    def test_normal(self):
        for d in self.docs:
            input_test_number = input('>>>')
            pass
        # self.assertEqual(1, d.fizzbuzz(1))



if __name__ == "__main__":
    unittest.main()