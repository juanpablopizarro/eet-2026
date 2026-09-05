# 🚀 3-Week Capstone Challenge: AI-Assisted Game Engineering

Welcome to your final engineering sprint. Over the next three weeks, your team of 2–3 will conceptualize, architect, and ship a fully functional 2D game.

The twist: **You are required to use Artificial Intelligence as your primary development co-pilot.**

You will act as software architects and lead engineers directing AI assistants to scaffold, generate assets, debug edge cases, and write automated tests.

However: **If you ship code you cannot explain, you fail.**

---

## 🎯 The Rules of Engagement

1. **Stack & Game Freedom:** Your team chooses the game concept (Space Invaders, Platformer, Top-Down Shooter, Puzzle, etc.) and the tech stack (Python, TypeScript, Godot/C#, Rust, etc.).
2. **Architectural Ownership:** Your team defines its own project scaffolding. Monolithic, single-file scripts are banned. Whatever layout you choose, it must honor four separation-of-concerns invariants:
   * **Documentation & Audit:** Clear architectural justification and an active AI prompt trail.
   * **Asset Isolation:** Separation of multimedia files from source code.
   * **Configuration Decoupling:** Global constants (speeds, screen bounds, keys) extracted from core logic.
   * **Automated Testing:** Dedicated test files decoupled from presentation/rendering loops.
3. **The Co-Pilot Mandate:** You must deliberately leverage AI across three distinct sprint phases:
   * **Week 1 (Design & Setup):** Brainstorming mechanics, generating architecture diagrams, establishing schemas, and synthesizing initial assets.
   * **Week 2 (Core Development):** Writing discrete algorithms, building entity systems, and using AI for targeted problem-solving.
   * **Week 3 (Testing & Bugfixing):** Generating unit test suites, profiling performance, and auditing edge-case bugs.
4. **Mandatory Audit Trail:** Every significant AI prompt, hallucination, and human refactor must be logged in `docs/AI_LOG.md`.
5. **The Oral Defense:** During the final demo, the teacher will select random functions from your repository. Any team member must be able to explain the logic line-by-line on the spot.
6. **Github:** Project must be delivered in a github.com repo including all the instructions to build, test and run the project.

---

## 🏗️ Project Scaffolding: Design It Your Way

You are free to design the directory layout that best fits your chosen engine, build tools, and team workflow. Document and defend your choice in `docs/ARCHITECTURE.md`.

Below are three example scaffolds reflecting standard industry conventions across different ecosystems:

### Example A: Python / Pygame (Classic Modular)
```text
space-defender-py/
├── docs/
│   ├── ARCHITECTURE.md
│   └── AI_LOG.md
├── tests/
│   ├── test_player_bounds.py
│   └── test_score_system.py
├── src/
│   ├── config.py
│   ├── entities/
│   │   ├── player.py
│   │   └── alien.py
│   ├── systems/
│   │   ├── collision.py
│   │   └── wave_spawner.py
│   └── main.py
├── assets/
│   ├── sprites/
│   └── audio/
├── requirements.txt
├── .gitignore
└── README.md
```

### Example B: TypeScript / Web Game (Component & Asset Pipeline)
```text
retro-invaders-web/
├── docs/
│   ├── ARCHITECTURE.md
│   └── AI_LOG.md
├── tests/
│   └── combat.test.ts
├── public/
│   └── assets/
│       ├── images/
│       └── sfx/
├── src/
│   ├── constants/
│   │   └── gameSettings.ts
│   ├── components/
│   │   ├── Player.ts
│   │   └── EnemyFleet.ts
│   ├── core/
│   │   ├── GameLoop.ts
│   │   └── InputManager.ts
│   └── index.ts
├── package.json
├── tsconfig.json
├── .gitignore
└── README.md
```
### Example C: Godot Engine / C# (Scene-Driven Pattern)
```text
pixel-rogue-godot/
├── docs/
│   ├── ARCHITECTURE.md
│   └── AI_LOG.md
├── Tests/
│   └── TestHealthSystem.cs
├── Scenes/
│   ├── UI/
│   │   └── HUD.tscn
│   ├── Characters/
│   │   ├── Player.tscn
│   │   └── Enemy.tscn
│   └── Levels/
│       └── Arena.tscn
├── Scripts/
│   ├── Config/
│   │   └── GameConstants.cs
│   ├── Systems/
│   │   └── HitboxHandler.cs
│   └── GameController.cs
├── Assets/
│   ├── Textures/
│   └── Sounds/
├── project.godot
├── .gitignore
└── README.md
```
