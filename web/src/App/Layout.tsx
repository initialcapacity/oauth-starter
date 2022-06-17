import * as React from 'react';

export const Layout = ({children}: { children: React.ReactElement }): React.ReactElement =>
    <section className={'App'}>
        <header>
            <div className="container">
                <h1>OAuth SPA starter</h1>
            </div>
        </header>
        <section className="callout">
            <div className="container">
                an <span className="branded">oauth[]</span> experiment
            </div>
        </section>
        <main>
            {children}
        </main>
        <footer>
            <div className="container">
                <script>document.write("©" + new Date().getFullYear());</script>
                Initial Capacity, Inc. All rights reserved.
            </div>
        </footer>
    </section>;
