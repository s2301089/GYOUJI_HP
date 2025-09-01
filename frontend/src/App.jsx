import React from 'react';
import './index.css';

const App = () => {
  const scrollToSection = (id) => {
    const element = document.getElementById(id);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth' });
    }
  };

  return (
    <div className="bg-gray-100 text-gray-800 font-sans">
      {/* Header */}
      <header className="bg-gray-800 text-white shadow-md">
        <div className="container mx-auto px-6 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold">行事委員会</h1>
          <nav className="hidden md:flex space-x-4">
            <button onClick={() => scrollToSection('home')} className="hover:text-gray-300">ホーム</button>
            <button onClick={() => scrollToSection('about')} className="hover:text-gray-300">概要</button>
            <button onClick={() => scrollToSection('events')} className="hover:text-gray-300">イベント</button>
            <button onClick={() => scrollToSection('roles')} className="hover:text-gray-300">役職</button>
            <button onClick={() => scrollToSection('join')} className="hover:text-gray-300">参加方法</button>
          </nav>
        </div>
      </header>

      {/* Hero Section */}
      <section id="home" className="bg-gray-700 text-white text-center py-20">
        <h1 className="text-4xl md:text-6xl font-bold">仙台高専広瀬キャンパス 行事委員会</h1>
        <p className="mt-4 text-lg md:text-xl">学校を、もっと面白く。君のアイデアで、最高のイベントを創ろう！</p>
      </section>

      {/* Main Content */}
      <main className="container mx-auto px-6 py-12">
        {/* About Section */}
        <section id="about" className="mb-12">
          <h2 className="text-3xl font-bold text-center mb-6">行事委員会へようこそ！</h2>
          <p className="text-lg text-center max-w-3xl mx-auto">
            私たち行事委員会は、学校生活を彩る様々なイベントを企画・運営しています。
            仲間と協力し、一つの目標に向かって全力で取り組む。その経験は、きっとあなたの大きな財産になります。
          </p>
        </section>

        {/* Events Section */}
        <section id="events" className="mb-12">
          <h2 className="text-3xl font-bold text-center mb-6">主なイベント</h2>
          <div className="grid md:grid-cols-2 gap-8">
            <div className="bg-white p-6 rounded-lg shadow-md">
              <h3 className="text-2xl font-bold mb-4">スポーツ大会</h3>
              <ul className="list-disc list-inside mb-4">
                <li>春季スポーツ大会</li>
                <li>秋季スポーツ大会</li>
              </ul>
              <h4 className="text-xl font-semibold mb-2">運営の流れ</h4>
              <ol className="list-decimal list-inside space-y-2">
                <li>行う競技を決定する(実現可能な競技を模索します)</li>
                <li>競技長を決める</li>
                <li>競技長が競技メンバーと相談しながら競技要項を作成する。</li>
                <li>物品リストを作成する</li>
                <li>競技メンバーのシフトを作成する</li>
                <li>前日に大会準備を進める</li>
                <li>大会当日の朝に最終的な準備をし開会式をする</li>
                <li>大会結果を集計し、閉会式で結果発表を行う</li>
                <li>後日の委員会で反省点を共有し、のちの大会運営に生かす</li>
              </ol>
            </div>
            <div className="bg-white p-6 rounded-lg shadow-md">
              <h3 className="text-2xl font-bold mb-4">その他のイベント</h3>
               <ul className="list-disc list-inside mb-4">
                <li>予餞会</li>
                <li>新入生歓迎会</li>
              </ul>
              <h4 className="text-xl font-semibold mb-2">予餞会運営の流れ</h4>
              <p  className="mb-4">予餞会では、各部活動や対象の教員からビデオレターを募りそれを流すという形式をとっている</p>
              <ol className="list-decimal list-inside space-y-2">
                <li>部活動と対象となる教員を整理する</li>
                <li>委員でグループを作り、それぞれいくつかの部活動等を担当する</li>
                <li>担当対象の部活動や教員にビデオレターをお願いする</li>
                <li>ビデオレターを回収し、まとめる</li>
                <li>予餞会当日にまとめてビデオレターを流す</li>
              </ol>
            </div>
          </div>
        </section>

        {/* Roles Section */}
        <section id="roles" className="mb-12">
          <h2 className="text-3xl font-bold text-center mb-6">役職紹介</h2>
          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8 text-center">
            <div className="bg-white p-6 rounded-lg shadow-md">
              <h3 className="text-2xl font-bold mb-2">行事委員長</h3>
              <p>委員会全体のリーダーとして、活動をまとめます。</p>
            </div>
            <div className="bg-white p-6 rounded-lg shadow-md">
              <h3 className="text-2xl font-bold mb-2">副行事委員長</h3>
              <p>委員長を補佐し、円滑な運営をサポートします。</p>
            </div>
            <div className="bg-white p-6 rounded-lg shadow-md">
              <h3 className="text-2xl font-bold mb-2">会計</h3>
              <p>委員会の予算を管理し、活動資金を確保します。</p>
            </div>
            <div className="bg-white p-6 rounded-lg shadow-md">
              <h3 className="text-2xl font-bold mb-2">書記</h3>
              <p>会議の議事録を作成し、活動の記録を残します。</p>
            </div>
          </div>
        </section>

        {/* Join Section */}
        <section id="join" className="text-center bg-white py-12 px-6 rounded-lg shadow-md">
          <h2 className="text-3xl font-bold mb-4">君も仲間にならないか？</h2>
          <p className="text-lg mb-2">行事委員会は、各クラスから2名まで参加できます。</p>
          <p className="text-lg mb-4">学校を盛り上げたい、新しいことに挑戦したい、そんな熱い想いを持った君を待っています！</p>
          <p className="text-lg">興味のある方は、来年に行事委員会に立候補してみてください。</p>
        </section>
      </main>

      {/* Footer */}
      <footer className="bg-gray-800 text-white text-center py-4 mt-12">
        <p>&copy; 2025 仙台高専広瀬キャンパス 行事委員会</p>
      </footer>
    </div>
  );
};

export default App;
