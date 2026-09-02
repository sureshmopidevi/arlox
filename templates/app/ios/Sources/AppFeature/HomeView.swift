import SwiftUI

public struct HomeView: View {
    @StateObject private var viewModel: HomeViewModel

    public init(viewModel: HomeViewModel = HomeViewModel()) {
        _viewModel = StateObject(wrappedValue: viewModel)
    }

    public var body: some View {
        VStack(spacing: 12) {
            Text(viewModel.greeting)
                .font(.largeTitle.bold())
                .multilineTextAlignment(.center)
            Text("SwiftUI · MVVM · Repository")
                .foregroundStyle(.secondary)
        }
        .padding()
        .task { await viewModel.load() }
    }
}
